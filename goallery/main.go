package main

import (
	"fmt"
	"net/http"

	"devisions.org/goallery/controllers"
	"github.com/gorilla/mux"
)

func main() {

	staticCtrl := controllers.NewStatic()
	usersCtrl := controllers.NewUsers()

	r := mux.NewRouter()

	r.Handle("/", staticCtrl.HomeView).Methods(http.MethodGet)
	r.Handle("/contact", staticCtrl.ContactView).Methods(http.MethodGet)

	r.HandleFunc("/signup", usersCtrl.New).Methods(http.MethodGet)
	r.HandleFunc("/signup", usersCtrl.Create).Methods(http.MethodPost)

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	_ = http.ListenAndServe(":3000", r)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, "<h3>Requested Page was not found</h3>")
}

// A helper function that panics on error.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
