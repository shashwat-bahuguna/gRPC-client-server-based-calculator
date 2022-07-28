package main

import (
	"context"
	"fmt"
	"log"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"gRPC-based-calculator/messages_proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Initiating Server")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := messages_proto.NewCalculatorServiceClient(cc)
	Sum(c)
}

func Sum(c messages_proto.CalculatorServiceClient) {
	fmt.Println("Unary GRPC, addition example")

	req := messages_proto.SumRequest{
		NumA: 10,
		NumB: 20,
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while sum call: %v", err)
	}

	log.Printf("Response from Server: %v", resp.Result)

}
