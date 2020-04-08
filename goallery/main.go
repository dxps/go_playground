package main

import (
	"fmt"
	"net/http"

	"devisions.org/goallery/controllers"
	"devisions.org/goallery/views"
	"github.com/gorilla/mux"
)

var homeView, contactView *views.View

func main() {

	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	usersCtrl := controllers.NewUsers()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/contact", contactHandler)
	r.HandleFunc("/signup", usersCtrl.New)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	_ = http.ListenAndServe(":3000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
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
