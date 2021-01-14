package main

import (
	"log"
	"net/http"

	"github.com/devisions/go-playground/go-chi-rate-limit/traffic"
	"github.com/go-chi/chi"
)

const ADDR = ":8000"

var trafficCtrl = traffic.NewIPRateControl(1, 1)

func main() {

	r := chi.NewRouter()

	r.Get("/", rootHandler)

	go trafficCtrl.PeriodicCleanup()

	log.Println(">>> Starting HTTP Server listening on", ADDR)
	if err := http.ListenAndServe(ADDR, r); err != nil {
		log.Fatalln(">>> Unable to start the HTTP Server. Reason:", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rl := trafficCtrl.GetLimiter(r.RemoteAddr)
	if !rl.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		log.Println(">>> Rejecting with", http.StatusTooManyRequests)
		return
	}
	_, _ = w.Write([]byte("root"))
}
