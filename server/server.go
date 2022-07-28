package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"gRPC-based-calculator/messages_proto"
)

type server struct {
	messages_proto.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *messages_proto.SumRequest) (resp *messages_proto.SumResponse, err error) {
	fmt.Println("Sum Request Recieved")

	resp = &messages_proto.SumResponse{}
	resp.Result = req.NumA + req.NumB
	err = nil
	return resp, err
}

func main() {
	fmt.Println("Initiating Server")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	fmt.Println("Successfully initiated Server")

	s := grpc.NewServer()
	messages_proto.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server
	// reflection.Register(s)

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
