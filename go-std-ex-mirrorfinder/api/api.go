package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"devisions.org/go-playground/go-std-ex-mirrorfinder/mirrors"
)

// Response is what gets returned to the client as JSON.
type Response struct {
	FastestURL string        `json:"fastest_url"`
	Latency    time.Duration `json:"latency"`
}

func findFastest(urls []string) Response {
	urlChan := make(chan string)
	latencyChan := make(chan time.Duration)

	for _, url := range urls {
		mirrorURL := url
		go func() {
			start := time.Now()
			_, err := http.Get(mirrorURL + "/README")
			latency := time.Now().Sub(start) / time.Millisecond
			if err == nil {
				urlChan <- mirrorURL
				latencyChan <- latency
			}
		}()
	}
	return Response{<-urlChan, <-latencyChan}
}

// FastestMirrorHandler is the request handler for /fastest-mirror.
func FastestMirrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(">>> Got request '%s%s'", r.Host, r.URL.RequestURI())
	response := findFastest(mirrors.MirrorList)
	respJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(respJSON)
	log.Printf(">>> Sent response '%s'", respJSON)
}
