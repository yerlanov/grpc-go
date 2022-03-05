package main

import (
	"context"
	"fmt"
	"github.com/yerlanov/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("CLIENT")

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}

	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(c)

	//doStreamServer(c)

	//doClientStream(c)

	//doBiDirectionalStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		Num1: 1,
		Num2: 2,
	}
	res, err := c.SumTwoNumbers(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling sum rpc: %v", err)
		return
	}

	log.Printf("Response from calcualator: %v", res.Res)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SquareRootRequest{
		Num: 10 ,
	}

	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			return
		}
		log.Fatalf("error while calling sum rpc: %v", err)
		return
	}

	log.Printf("Response from calcualator: %v", res.Res)
}

func doStreamServer(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberRequest{Num: 1246123441241243720}

	res, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling rpc: %v", err)
	}

	for {
		msg, err := res.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("error while reading stream: %v: ", err)
		}
		fmt.Printf("Response: %v\n", msg.Res)
	}
}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.ComputeAverageRequest{
		{Num: 1},
		{Num: 2},
		{Num: 3},
		{Num: 4},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling rpc: %v", err)
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

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.FindMaximumRequest{
		{Num: 1},
		{Num: 5},
		{Num: 3},
		{Num: 6},
		{Num: 2},
		{Num: 20},
	}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling rpc: %v", err)
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
				}
			}

			fmt.Printf("received: %v\n", res.Res)
		}
		close(waitCh)
	}()

	<-waitCh
}
