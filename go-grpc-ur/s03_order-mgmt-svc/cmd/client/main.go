package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/devisions/go-playground/go-grpc-ur/s03-order-mgmt-svc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	address = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		abort("Failed to connect to Server. Reason: %s", err)
	}
	defer conn.Close()

	client := pb.NewOrderMgmtClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get an order.
	ordId := "100"
	order, err := client.GetOrder(ctx, &wrapperspb.StringValue{Value: ordId})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.NotFound {
				abort(fmt.Sprintf("Failed to get order. Reason: %s", err))
			}
			log.Printf(">>> Order with id '%s' was not found.", ordId)
		}
	} else {
		fmt.Printf("Got order: %+v", order)
	}
}

func abort(msg string, data ...interface{}) {
	log.Printf(msg, data...)
	os.Exit(1)
}
