package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	calculatorpb "RPC/calculator/calculatorpb"

	_ "golang.org/x/tools/go/analysis/passes/nilfunc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
}

func (*Server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Recieved Calulator Rpc : %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil

}
func (*Server) PrimeNumberDecompostion(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompostionServer) error {
	fmt.Printf("Recieved Prime Number Decomposition factor: %v", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Println("Divisor Increamented")
		}
	}
	return nil

}
func (*Server) ComputeAverge(stream calculatorpb.CalculatorService_ComputeAvergeServer) error {
	fmt.Println("Inside ComputeAverage Pb")
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error While Reading Client Stream:%v", err)
			return nil
		}
		sum += req.GetNumber()
		count++

	}

}
func (*Server) FindMaxminum(stream calculatorpb.CalculatorService_FindMaxminumServer) error {
	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error While Reading Client Stream:%v", err)
			return err
		}
		if maximum < req.GetNumber() {
			maximum = req.GetNumber()
		}
		// fmt.Printf("New Number:%v\n", maximum)
		sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: maximum,
		})
		if sendErr != nil {
			log.Fatalf("Error While Sending Data to Client :%v", sendErr)
			return sendErr
		}

	}

}
func (*Server) SquareRoot(ctx context.Context, req *calculatorpb.SquarerootRequest) (*calculatorpb.SquarerootResponse, error) {
	fmt.Println("Recieved Square Root Rpc")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Recieved a Negative Number: %v\n", number))
	} else {
		return &calculatorpb.SquarerootResponse{NumberRoot: math.Sqrt(float64(number))}, nil
	}

}
func main() {
	fmt.Println("Inside Calculatorpb")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Unable to Start Server : %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Error Listing to a Server : %v", err)
	}

}
