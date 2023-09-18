package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main(){
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to start server")
	}
	grpcServer := grpc.NewServer()
	log.Printf("Server running at 9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v",err)
	}
}