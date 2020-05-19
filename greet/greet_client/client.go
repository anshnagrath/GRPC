package main

import (
	"RPC/greet/greetpb"

	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	// to tell grpc to not to use ssl
	certfile := "./ca.crt"
	creds, sslErr := credentials.NewClientTLSFromFile(certfile, "")
	if sslErr != nil {
		log.Fatalf("SSL Error:%v", sslErr)
		return
	} else {
		clientconn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("Could not connect : %v", err)
		}
		c := greetpb.NewGreetingServiceClient(clientconn)
		defer clientconn.Close()

		doUrinary(c)
		// doServerStreaming(c)
		// doClientStreaming(c)
		// doBiDiStreaming(c)
		// doUniaryWithTimeout(c, 1)
		// doUniaryWithTimeout(c, 3)
		// doUniaryWithTimeout(c, 7)
	}

}
func doUrinary(c greetpb.GreetingServiceClient) {
	fmt.Printf("Starting to do a Unary RPC:%v", c)
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

func doServerStreaming(c greetpb.GreetingServiceClient) {
	req := &greetpb.GreetManyRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ansh",
			LastName:  "Nagrath",
		},
	}
	stream, err := c.GreetMany(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Calling GreetMany Service %v", err)
	}
	log.Printf("Response from GreetMany : %v", stream)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("Stream End %v", err)
		}
		log.Printf("Response From Greet Many Times:%v", msg.GetResult())
	}

}
func doClientStreaming(c greetpb.GreetingServiceClient) {
	stream, err := c.LongGreet(context.Background())
	req := []*greetpb.LongGreetManyRequest{
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "YO",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "stark",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "what up",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "having fun ?",
			},
		},
	}
	if err != nil {
		log.Fatalf("Error while fetching long greet:%v", err)
	}
	for _, r := range req {
		stream.Send(r)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving response from long greet : %v", err)
	}
	fmt.Printf("Long Greet Response : %v\n", res)

}

func doBiDiStreaming(c greetpb.GreetingServiceClient) {
	fmt.Println("Started Streaming Data")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error While Calling GreetMany Service %v", err)
		return
	}
	waitc := make(chan struct{})

	go func() {
		//function to send a bunch of messages
		req := []*greetpb.GreeteveryoneRequest{
			&greetpb.GreeteveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "YO",
				},
			},
			&greetpb.GreeteveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "stark",
				},
			},
			&greetpb.GreeteveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "what up",
				},
			},
			&greetpb.GreeteveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "having fun ?",
				},
			},
		}
		if err != nil {
			log.Fatalf("Error while fetching long greet:%v", err)
		}
		for _, r := range req {
			stream.Send(r)
			fmt.Println("Sending Message : \n", r)
		}
		stream.CloseSend()

	}()

	go func() {
		// function to recive a bunch of message
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break

			}
			if err != nil {
				log.Fatalf("Error while fetching long greet:%v", err)
				break
			}
			fmt.Printf("Recieved: %v\n", res.GetResult())
		}
		close(waitc)

	}()
	<-waitc
}

func doUniaryWithTimeout(c greetpb.GreetingServiceClient, seconds time.Duration) {
	fmt.Println("Trying to make a Uniary With Deadline")
	ctx, cancel := context.WithTimeout(context.Background(), seconds*time.Second)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "YO",
		},
	})
	if err != nil {
		resErr, _ := status.FromError(err)
		if resErr.Code() == codes.DeadlineExceeded {
			log.Fatalf("Error While Making An Timeout RPC Status Code:%v", resErr.Code())
			log.Fatalf("Error While Making An Timeout RPC:%v", resErr.Message())
			return
		} else {
			log.Fatalf("Error While Making An Timeout RPC Status Code:%v", resErr.Code())
			log.Fatalf("Error While Making An Timeout RPC:%v", resErr.Message())
			return
		}
	} else {
		fmt.Printf("Response Recieved :%v", res.Result)

	}
}
