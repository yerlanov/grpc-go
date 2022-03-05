package main

import (
	"context"
	"fmt"
	"github.com/yerlanov/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("CLIENT")

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDirectionalStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Raim",
			LastName:  "some",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greed rpc: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Raim",
			LastName:  "some",
		},
	}
	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greed rpc: %v", err)
	}

	for {
		msg, err := res.Recv()
		if err != nil {
			if err == io.EOF {
				//
				break
			}
			log.Fatalf("error while reading stream: %v", err)
		}
		fmt.Printf("Response: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Raim",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dina",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Baha",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Almas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bola",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling long greet: %v", err)
		return
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v", req)
		err = stream.Send(req)
		if err != nil {
			log.Fatalf("error while sending data in stream: %v", err)
			return
		}
		time.Sleep(1 * time.Second)
	}

	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
		return
	}
	fmt.Printf("Response: %v\n", recv)

}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Raim",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dina",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Baha",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Almas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bola",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling greed rpc: %v", err)
		return
	}

	waitCh := make(chan struct{})

	go func() {
		for _, request := range requests {
			fmt.Printf("Sending: %v\n", request)
			err = stream.Send(request)
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
		err = stream.CloseSend()
		if err != nil {
			return
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatalf("error while receiving: %v", err)
					break
				}
			}

			fmt.Printf("received: %v\n", res.GetResult())
		}
		close(waitCh)
	}()

	<-waitCh
}
