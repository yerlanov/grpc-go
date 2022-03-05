package main

import (
	"context"
	"fmt"
	"github.com/yerlanov/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct{}

func (s *server) SquareRoot(ctx context.Context, request *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function was invoked")
	number := request.GetNum()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received negative number: %v", number))
	}

	return &calculatorpb.SquareRootResponse{Res: math.Sqrt(float64(number))}, nil
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("GreetEveryone function was invoked")
	maxNum := int64(0)
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
		current := req.Num
		if current < maxNum {
			continue
		}
		err = stream.Send(&calculatorpb.FindMaximumResponse{Res: current})
		if err != nil {
			log.Fatalf("error while sending data to client: %v", err)
			return err
		}
		maxNum = current
	}
}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("LongGreet function was invoked")
	result := int64(0)
	counter := int64(0)
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				res := float64(result) / float64(counter)
				return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{Res: res})
			} else {
				log.Fatalf("error while reading client stream: %v", err)
			}
		}
		result += req.Num
		counter++
	}
}

func (s *server) SumTwoNumbers(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	//lastName := req.GetGreeting().GetLastName()
	result := req.GetNum1() + req.GetNum2()
	res := &calculatorpb.SumResponse{Res: result}
	return res, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	var (
		n = req.GetNum()
		k = int64(2)
	)
	for n > 1 {
		if n%k == 0 {
			err := stream.Send(&calculatorpb.StreamPrimeNumberResponse{Res: k})
			if err != nil {
				return err
			}
			n = n / k
		} else {
			k = k + 1
			fmt.Printf("Devisor has increased to: %v\n", k)
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
