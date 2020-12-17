package main

import (
	"log"
	"net/http"

	"github.com/devisions/go-playground/go-expvar-metrics/metrics"
)

func main() {

	http.HandleFunc("/", handleAll)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAll(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("..."))
	metrics.IncreaseRequests()
}
