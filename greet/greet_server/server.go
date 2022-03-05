package main

import (
	"context"
	"fmt"
	"github.com/yerlanov/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				log.Fatalf("error while reading client stream: %v", err)
				return err
			}
		}
		firstName := req.Greeting.GetFirstName()
		result := "Hello " + firstName + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{Result: result})
		if err != nil {
			log.Fatalf("error while sending data to client: %v", err)
			return err
		}
	}
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{Result: result}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	req.Greeting.GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{Result: fmt.Sprintf("%d: Hello %s", i, firstName)}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked")
	result := ""
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
			} else {
				log.Fatalf("error while reading client stream: %v", err)
			}
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
