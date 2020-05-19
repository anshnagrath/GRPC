package main

import (
	greetpb "RPC/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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

func (*server) GreetMany(req *greetpb.GreetManyRequest, stream greetpb.GreetingService_GreetManyServer) error {
	fmt.Printf("Greet Many Times Envoked: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyResponse{
			Result: "Hello " + firstName + " number " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(time.Second * 1)
	}
	return nil
}
func (*server) LongGreet(stream greetpb.GreetingService_LongGreetServer) error {
	result := " "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetManyResponse{
				Result: result,
			})
		} else if err != nil {
			log.Fatalf("File Ended:%v", err)
		} else {
			firstName := req.GetGreeting().GetFirstName()
			result += " Hello " + firstName + "!"
		}

	}
}
func (*server) GreetEveryone(stream greetpb.GreetingService_GreetEveryoneServer) error {
	fmt.Println("Invocked GreetEveryOne Request")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error While Recieving Stream From Client: %v", err)
			return err
		}
		firstName := req.GetGreeting().FirstName
		result := "Hello " + firstName + " !"
		senderr := stream.Send(&greetpb.GreeteveryoneResponse{
			Result: result,
		})
		if senderr != nil {
			log.Fatalf("Error While Recieving Stream From Client: %v", senderr)
			return senderr
		}

	}

}
func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreatWithDeadlineResponse, error) {
	fmt.Println("Greet Response with TimeOut")
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 3; i++ {

		if ctx.Err() == context.Canceled {
			fmt.Println("Client Canceled the Request!")
			return nil, status.Error(codes.DeadlineExceeded, "Error Client Canceled the Request")
		}
		time.Sleep(1 * time.Second)

	}
	Greeting := "Hello " + firstName
	response := &greetpb.GreatWithDeadlineResponse{
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
	certFile := "server.crt"
	keyFile := "server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Sll Error :%v", sslErr)
		return
	} else {

		s := grpc.NewServer(grpc.Creds(creds))
		greetpb.RegisterGreetingServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		} else {

		}
	}

}
