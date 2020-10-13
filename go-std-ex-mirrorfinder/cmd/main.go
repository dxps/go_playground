//
// main logic of the API server
//
package main

import (
	"log"
	"net/http"
	"time"

	"devisions.org/go-playground/go-std-ex-mirrorfinder/api"
)

func main() {

	http.HandleFunc("/fastest-mirror", api.FastestMirrorHandler)
	port := ":8001"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf(">>> Starting listening on port %s\n", port)
	// Start the server and log as fatal a possible error
	// that may be returned due to an unsuccessful launch.
	log.Fatal(server.ListenAndServe())
}
