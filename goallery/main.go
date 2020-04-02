package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h3>Goallery begins!</h3>")
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":3000", nil)
}
