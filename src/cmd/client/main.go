package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"team00/client"
	frequency "team00/generated"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DefaultPostgresHost     string = "localhost"
	DefaultPostgresUser     string = "postgres"
	DefaultPostgresPassword string = "postgres"
	DefaultPostgresDatabase string = "postgres"
	DefaultPostgresPort     int    = 5432
)

func main() {
	args := arguments{
		Address: flag.String("address", "0.0.0.0:8888", "gRPC server address in form address:port"),
		K:       flag.Float64("k", 1, "Multiplier for STD"),
	}
	flag.Parse()

	conn, err := grpc.NewClient(
		*args.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := frequency.NewFrequencyRandomizerClient(conn)

	stream, err := c.SpawnFrequencies(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	dbCredentials := getDatabaseCredentialsFromEnv()
	log.Printf(
		"postgres info: host=%s user=%s database=%s port=%d\n",
		dbCredentials.Host,
		dbCredentials.User,
		dbCredentials.DatabaseName,
		dbCredentials.Port,
	)

	dbConn, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		dbCredentials.Host,
		dbCredentials.User,
		dbCredentials.Password,
		dbCredentials.DatabaseName,
		dbCredentials.Port,
	)), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	dbConn.AutoMigrate(&frequency.MessageORM{})

	processor := client.NewClient()
	err = processor.Calibrate(stream, client.DetectSamplesOptimalSize)
	if err != nil {
		log.Fatalln(errors.Join(fmt.Errorf("failed to calibrate"), err))
	}

	for {
		select {
		case <-stream.Context().Done():
			log.Println("server is closed")
		default:
			message, err := stream.Recv()
			if err != nil {
				log.Fatalln(err)
			}

			isAnomaly, diff, err := processor.DetectAnomaly(message.GetFrequency(), *args.K)
			if err != nil {
				log.Fatalln(err)
			}

			if isAnomaly {
				fmt.Printf(
					"anomaly detected: frequency=%f, diff=%f, k=%f, mean=%f, std=%f\n",
					message.GetFrequency(),
					diff,
					*args.K,
					processor.Detector.Mean(),
					processor.Detector.Std(),
				)

				ormMessage, err := message.ToORM(context.Background())
				if err != nil {
					log.Fatalln(err)
				}
				dbConn.Create(&ormMessage)
			}
		}
	}
}

type DbConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         int
}

func getDatabaseCredentialsFromEnv() DbConfig {
	result := DbConfig{
		Host:         os.Getenv("POSTGRES_HOST"),
		User:         os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("POSTGRES_DATABASE"),
	}
	port, err := strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 32)
	if err != nil {
		port = int64(DefaultPostgresPort)
	}
	result.Port = int(port)

	if result.Host == "" {
		result.Host = DefaultPostgresHost
	}
	if result.User == "" {
		result.User = DefaultPostgresUser
	}
	if result.Password == "" {
		result.Password = DefaultPostgresPassword
	}
	if result.DatabaseName == "" {
		result.DatabaseName = DefaultPostgresDatabase
	}

	return result
}

type arguments struct {
	Address *string
	K       *float64
}
