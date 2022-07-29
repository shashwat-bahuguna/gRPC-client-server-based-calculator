package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"gRPC-based-calculator/messages_proto"
)

const INT32_MIN int32 = -2147483648

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

func isPrime(num int32) bool {
	var i int32
	for i = 2; i < num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func (*server) PrimeLister(req *messages_proto.PrimeRequest, resp_interface messages_proto.CalculatorService_PrimeListerServer) error {
	fmt.Println("Prime Listing Request Recieved")

	num := req.Num
	var i int32
	resp := messages_proto.PrimeResponse{}
	for i = 2; i <= num; i++ {
		if num%i == 0 && isPrime(i) {
			resp.Result = i
			resp_interface.Send(&resp)
			fmt.Printf("Sending to Client: %v\n", resp.Result)
		}
	}

	fmt.Printf("Prime Listing Request Completed\n")
	return nil
}

func (*server) AverageCalculator(stream messages_proto.CalculatorService_AverageCalculatorServer) error {
	fmt.Println("Average Calculator Request Recieved")

	var n int32 = 0
	var sum int32 = 0
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading client stream
			fmt.Printf("Average Calculator Request Completed, Res: %v\n", sum/n)
			return stream.SendAndClose(&messages_proto.AverageResponse{Result: sum / n})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream : %v\n", err)
			return err
		}
		n++
		sum += msg.Num
	}
}

func (*server) MaxCalculator(stream messages_proto.CalculatorService_MaxCalculatorServer) error {
	fmt.Println("Max Calculator Request Recieved")
	defer fmt.Println("Max Calculator Request Completed")

	resp := messages_proto.MaxResponse{}
	resp.Result = INT32_MIN
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while recieving message in max calculator: %v\n", err)
			return err
		}

		if msg.GetNum() > resp.Result {
			resp.Result = msg.GetNum()
		}

		fmt.Printf("Writing to client: %v\n", resp.Result)
		err = stream.Send(&resp)
		if err != nil {
			log.Fatalf("Error while sending response in Send: %v\n", err)
			return err
		}

	}
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
