package main

import (
	"log"
	"os"

	"github.com/devisions/go-playground/go-grpc-ur/s03-order-mgmt-svc/cmd/server/internal"
)

const (
	address = ":50051"
)

func main() {

	_, err := internal.StartServer(address)
	if err != nil {
		abort(err.Error())
	}
}

func abort(msg string, data ...interface{}) {
	log.Printf(msg, data...)
	os.Exit(1)
}
