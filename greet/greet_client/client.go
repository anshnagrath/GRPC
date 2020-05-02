package main

import (
	"RPC/greet/greetpb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am in")
	// to tell grpc to not to use ssl
	clientconn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	c := greetpb.NewGreetingServiceClient(clientconn)
	defer clientconn.Close()
	fmt.Printf("Created Client:%v", c)
	doUrinary(c)

}
func doUrinary(c greetpb.GreetingServiceClient) {
	fmt.Println("Starting to do a Unary RPC:%v", c)
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ansh",
			LastName:  "Nagrath",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Calling Greet %v", err)
	}
	log.Printf("Response from Greet : %v", res)

}
