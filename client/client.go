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
	Sum_tester(c)
	PrimeLister_Tester(c)
	AverageCalculator_Tester(c)
	MaxCaculator_Tester(c)
}

func Sum_tester(c messages_proto.CalculatorServiceClient) {
	fmt.Printf("\n\t\t#####################   SUM TEST (UNARY) ##########################\n")
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
	fmt.Printf("\n\t\t#####################   PRIME LISTING TEST (SERVER STREAMING) ##########################\n")
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

func AverageCalculator_Tester(c messages_proto.CalculatorServiceClient) {
	fmt.Printf("\n\t\t#####################   AVERAGE CALCULATOR TEST (CLIENT STREAMING) ##########################\n")
	data := [][]int32{{10, 20, 30, 10, 50}, {23, 74, 23, 48, 57}, {38, 50, 49, 99}}

	for _, slice := range data {
		AverageCalculator(c, slice)
	}
}

func AverageCalculator(c messages_proto.CalculatorServiceClient, data []int32) {
	log.Printf("Client Sreaming gRPC, Average Calculation, Target: %v", data)

	reqStream, err := c.AverageCalculator(context.Background())

	if err != nil {
		log.Fatalf("Error while initiating connection in Average Calculator: %v", err)
	}

	req := messages_proto.AverageRequest{}

	for _, val := range data {
		req.Num = val
		err = reqStream.Send(&req)
		if err != nil {
			log.Fatalf("Error while Sending in Average Calculator: %v", err)
		}
	}

	resp, err := reqStream.CloseAndRecv()
	log.Printf("Response from Client: %v", resp.GetResult())
}

func MaxCaculator_Tester(c messages_proto.CalculatorServiceClient) {
	fmt.Printf("\n\t\t#####################   MAX CALCULATOR TEST (BIDIRECTIONAL STREAMING) ##########################\n")
	data := [][]int32{{12, 43, 23, 54, 40, 60, 51, 70, 81}, {12, 12, 14, 29, 49, 30, 20, 29, 41, 32, 49}, {12, 23, 20, 40, 23, 11}, {23, 43, 22, 35, 54}}

	for _, slice := range data {
		MaxCaculator(c, slice)
	}
}

func MaxCaculator(c messages_proto.CalculatorServiceClient, data []int32) {
	log.Printf("Client Sreaming gRPC, Maximum Calculation, Target: %v", data)

	stream, err := c.MaxCalculator(context.TODO())

	if err != nil {
		log.Printf("Error while initiating reponse in Max Calculation: %v", err)
	}

	// Channel for blocking main thread
	ch := make(chan int)

	go func() {
		req := messages_proto.MaxRequest{}
		for _, val := range data {
			req.Num = val

			err = stream.Send(&req)

			if err != nil {
				log.Fatalf("Error while Sending in Max Calculator: %v", err)
			}
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				ch <- 1
				return
			}

			if err != nil {
				log.Printf("Error while receiving response in Max Calc: %v", err)
			}
			log.Printf("SERVER: %v ", msg.GetResult())
		}
	}()

	<-ch
}
