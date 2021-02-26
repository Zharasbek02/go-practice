package server

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*Server) PrimeDecomposition(req *calculatorpb.DecompositionRequest, stream calculatorpb.CalculatorService_DecompositionServer) error {
	fmt.Printf("PrimeDecomposition function was invoked with %v \n", req)
	number := req.GetNumber()
	var factor int64 = 2

	for number > 1 {
		if number%factor == 0 {
			res := &calculatorpb.DecompositionResponse{Result: factor}
			err := stream.Send(res)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
			number = number / factor
			time.Sleep(time.Second)
		} else {
			factor++
		}
	}

	return nil
}

func (*Server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Printf("Average function was invoked with a streaming request\n")

	var sum int64 = 0
	var count int64 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
