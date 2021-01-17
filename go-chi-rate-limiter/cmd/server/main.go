package main

import (
	"log"
	"net"
	"net/http"

	"github.com/devisions/go-playground/go-chi-rate-limit/traffic"
	"github.com/go-chi/chi"
)

const ADDR = ":8001"

var trafficCtrl = traffic.NewIPRateControl(0.5, 1)

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
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "internal issue", http.StatusInternalServerError)
	}
	rl := trafficCtrl.GetLimiter(ip)
	if !rl.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		log.Println(">>> Rejecting with", http.StatusTooManyRequests)
		return
	}
	_, _ = w.Write([]byte("root"))
}
