package main

import (
	"fmt"
	"net/http"

	"devisions.org/goallery/views"
	"github.com/gorilla/mux"
)

var homeView, contactView *views.View

func main() {

	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/contact", contactHandler)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	_ = http.ListenAndServe(":3000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil); err != nil {
		panic(err)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactView.Template.ExecuteTemplate(w, contactView.Layout, nil); err != nil {
		panic(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, "<h3>Requested Page was not found</h3>")
}
