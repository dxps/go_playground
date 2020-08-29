package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", articlePathHandler)
	r.HandleFunc("/articles", articlesQueryHandler)
	fmt.Println(">>> Starting the server ...")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// articlePathHandler is a Path-based handler.
// Request example: `curl localhost:8000/articles/books/123`
func articlePathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
	fmt.Fprintf(w, "ID: %v\n", vars["id"])
}

// articlesQueryHandler is a Query-based handler.
// Request example: `curl "localhost:8000/articles?id=123&category=books"`
func articlesQueryHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Got param id: %s\n", queryParams["id"][0])
	fmt.Fprintf(w, "Got param category: %s\n", queryParams["category"][0])
}
