//
// main logic of the API server
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"devisions.org/go-playground/go-std-ex-mirrorfinder/api"
)

func main() {

	http.HandleFunc("/fastest-mirror", api.FastestMirrorHandler)
	port := ":8000"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf(">>> Starting listening on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
