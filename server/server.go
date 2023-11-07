package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"places/places"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"google.golang.org/grpc"
)



func main(){
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to start server")
	}
	grpcServer := grpc.NewServer()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://username:anik32069@sa-project.10ptaex.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil{
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10 *time.Second)
	err = client.Connect(ctx)
	log.Printf("Connected to mongoDB")


	defer client.Disconnect(ctx)


	err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }
    databases, err := client.ListDatabaseNames(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
	placesDatabase := client.Database("places")
	placesCollection := placesDatabase.Collection("places")

	con := places.Connection{
		PlacesCollection: placesCollection,
	}
    fmt.Println(databases)
	fmt.Println(placesCollection)



	places.RegisterPlaceServiceServer(grpcServer, places.NewPlacesServer(con))
	log.Printf("Server running at 9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v",err)
	}
}