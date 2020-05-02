package main

import (
	calcculatorpb "RPC/calculator/calculatorpb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error While Connecting to Server : %v", err)
	}
	defer cc.Close()
	c := calcculatorpb.NewCalculatorServiceClient(cc)
	doUniary(c)

}
func doUniary(c calcculatorpb.CalculatorServiceClient) {
	fmt.Println("Inside Calculator unary")
	req := &calcculatorpb.SumRequest{
		FirstNumber:  23,
		SecondNumber: 345,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Connecting to Server : %v", err)
	}
	fmt.Println("Response from Cal RPC:", res.SumResult)

}
