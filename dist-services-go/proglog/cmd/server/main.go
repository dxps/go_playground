package main

import (
	"log"

	"github.com/devisions/go-playground/dist-services-go/proglog/internal/server"
)

func main() {
	log.Println(">>> Starting the HTTP Server and listening on port 8080 ...")
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
