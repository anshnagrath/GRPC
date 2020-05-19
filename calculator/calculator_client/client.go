package main

import (
	calcculatorpb "RPC/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error While Connecting to Server : %v", err)
	}
	defer cc.Close()
	c := calcculatorpb.NewCalculatorServiceClient(cc)
	// doUniary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doErrorUniary(c)

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
func doServerStreaming(c calcculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to a Server Streaming Request")
	req := &calcculatorpb.PrimeNumberDecompositionRequest{
		Number: 123,
	}
	stream, err := c.PrimeNumberDecompostion(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Creating Number:%v", err)
	} else {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Fatalf("Stream Ended:%v", msg)
				break
			} else {
				fmt.Println(msg.GetPrimeFactor())
			}

		}
	}

}
func doClientStreaming(c calcculatorpb.CalculatorServiceClient) {

	stream, err := c.ComputeAverge(context.Background())
	if err != nil {
		log.Fatalf("Error While Opening Stream:%v", err)
	}
	numbers := []int32{4, 5, 6, 34, 13, 4546, 78}

	for _, num := range numbers {
		stream.Send(&calcculatorpb.ComputeAverageRequest{
			Number: num,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error While Receving Response :%v", err)
	} else {
		fmt.Printf("Average:%v\n", res.GetAverage())
	}
}
func doBiDiStreaming(c calcculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaxminum(context.Background())
	if err != nil {
		log.Fatalf("Error Occured :%v", err)
	}

	waitC := make(chan struct{})

	go func() {
		//send multiple files
		numbers := []int32{4, 1, 434, 4, 5, 7, 2, 18, 19, 6, 32}
		for _, num := range numbers {
			// fmt.Println(i, num, "sccscsdc===")
			stream.Send(&calcculatorpb.FindMaximumRequest{Number: int32(num)})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		//reciee multiple files
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("All Results are here", res)
				break
			}
			if err != nil {
				log.Fatalf("Error While Receving Response :%v", err)
				break
			}
			max := res.GetMaximum()
			fmt.Printf("New Max:%v\n", max)

		}
		close(waitC)

	}()
	<-waitC

}

func doErrorUniary(c calcculatorpb.CalculatorServiceClient) {
	fmt.Println("Started to an Square Root Error Unary RPC ..")
	doSquareRoot(c, 12)
	doSquareRoot(c, -1)

}

func doSquareRoot(c calcculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calcculatorpb.SquarerootRequest{Number: number})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			//Actual thrown Error;
			fmt.Println("Error Status Code:\n", resErr.Code())
			fmt.Println("\n", resErr.Message())
			return

		} else {
			//framework  Type Error
			log.Fatalf("Big Error Calling a Sqaure Root:%v", err)
			return
		}

	}
	fmt.Printf("Square Root Recieved:%v", res.GetNumberRoot())
}
