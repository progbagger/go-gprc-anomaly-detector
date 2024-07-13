package main

import (
	"flag"
	"log"
	"net"
	frequency "team00/generated"
	"team00/server"

	"google.golang.org/grpc"
)

func main() {
	args := arguments{
		Address: flag.String("address", "0.0.0.0:8888", "gRPC server host in form address:port"),
	}
	flag.Parse()

	lis, err := net.Listen("tcp", *args.Address)
	if err != nil {
		log.Fatalln(err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	frequency.RegisterFrequencyRandomizerServer(
		grpcServer,
		&server.Server{},
	)

	log.Printf("starting gRPC server on %s", *args.Address)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}

type arguments struct {
	Address *string
}
