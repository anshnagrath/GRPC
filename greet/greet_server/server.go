package main

import (
	greetpb "RPC/greet/greetpb"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	Greeting := "Hello " + firstName
	response := &greetpb.GreetResponse{
		Result: Greeting,
	}
	return response, nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	fmt.Println("Server is up and Running at 50051")
	if err != nil {
		log.Fatalf("Failed To listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetingServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	} else {

	}

}
