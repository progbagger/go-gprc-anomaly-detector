package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"team00/client"
	frequency "team00/generated"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
			}
		}
	}
}

type arguments struct {
	Address *string
	K       *float64
}
