package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func doManyTimesFromServer(c calculatorpb.CalculatorServiceClient) {
	ctx := context.Background()
	req := &calculatorpb.PrimeDecompositionRequest{Number: 120}

	stream, err := c.Decomposition(ctx, req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break LOOP
		}
		if err != nil {
			log.Fatalf("error %v", err)
		}
		log.Printf("response:%v \n", res.GetResult())
	}
}

func doLongCalculateAverage(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.AverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	ctx := context.Background()
	stream, err := c.Average(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Average Response: %v\n", res)
}

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	doLongCalculateAverage(c)
}
