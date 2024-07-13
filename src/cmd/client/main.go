package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	frequency "team00/generated"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	args := arguments{
		Address: flag.String("address", "0.0.0.0:8888", "gRPC server address in form address:port"),
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

	client := frequency.NewFrequencyRandomizerClient(conn)

	response, err := client.SpawnFrequencies(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		message, err := response.Recv()
		if err != nil {
			log.Fatalln(err)
		}

		b, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
	}
}

type arguments struct {
	Address *string
}
