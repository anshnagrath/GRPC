package main

import (
	"context"
	"fmt"
	"log"
	"net"

	calculatorpb "RPC/calculator/calculatorpb"

	"google.golang.org/grpc"
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
