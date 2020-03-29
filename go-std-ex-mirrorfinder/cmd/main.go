//
// main logic of the API server
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"devisions.org/go-playground/go-std-ex-mirrorfinder/api"
	"devisions.org/go-playground/go-std-ex-mirrorfinder/mirrors"
)

func main() {
	http.HandleFunc("/fastest-mirror",
		func(w http.ResponseWriter, r *http.Request) {
			response := api.FindFastest(mirrors.MirrorList)
			respJson, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(respJson)
		})
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
