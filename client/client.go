package main

import (
	"context"
	"fmt"
	"io"
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
	// Sum_tester(c)
	// PrimeLister_Tester(c)
	MaxCalculator_Tester(c)
}

func Sum_tester(c messages_proto.CalculatorServiceClient) {
	data := [][]int32{{10, 20}, {20, 40}, {10, 15}, {100, 600}}
	for _, slice := range data {
		Sum(c, slice[0], slice[1])
	}
}

func Sum(c messages_proto.CalculatorServiceClient, a int32, b int32) {
	log.Printf("Unary gRPC, addition example, vals: %v, %v", a, b)

	req := messages_proto.SumRequest{
		NumA: a,
		NumB: b,
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while sum call: %v", err)
	}

	log.Printf("Response from Server: %v", resp.GetResult())

}

func PrimeLister_Tester(c messages_proto.CalculatorServiceClient) {
	data := []int32{210, 200, 100, 150}
	for _, val := range data {
		PrimeLister(c, val)
	}
}

func PrimeLister(c messages_proto.CalculatorServiceClient, num int32) {
	log.Printf("Server Sreaming gRPC, Prime Factors Calculation, Target: %v", num)

	req := messages_proto.PrimeRequest{
		Num: num,
	}

	respStream, err := c.PrimeLister(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while Prime call: %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("%v, ", msg.GetResult())
		// log.Printf("Response from Server: %v", resp.)
	}
}

func MaxCalculator_Tester(c messages_proto.CalculatorServiceClient) {
	data := [][]int32{{10, 20, 30, 10, 50}, {23, 74, 23, 48, 57}, {38, 50, 49, 99}}

	for _, slice := range data {
		MaxCalculator(c, slice)
	}
}

func MaxCalculator(c messages_proto.CalculatorServiceClient, data []int32) {

	log.Printf("Client Sreaming gRPC, Maximum Calculation, Target: %v", data)

	reqStream, err := c.MaxCalculator(context.Background())

	if err != nil {
		log.Fatalf("Error while initiating connection in Max Calculator: %v", err)
	}

	req := messages_proto.MaxRequest{}

	for _, val := range data {
		req.Num = val
		err = reqStream.Send(&req)
		if err != nil {
			log.Fatalf("Error while Sending in Max Calculator: %v", err)
		}
	}
	reqStream.CloseSend()

	resp, err := reqStream.Recv()
	log.Printf("Response from Client: %v", resp.GetResult())
}
