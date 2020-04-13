package main

import (
	"devisions.org/goallery/models"
	"fmt"
	"net/http"

	"devisions.org/goallery/controllers"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 54321
	user     = "goallery"
	password = "goallery"
	dbname   = "goallery"
)

func main() {

	// Database Init
	dbConnInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	userRepo, err := models.NewUserRepo(dbConnInfo)
	if err != nil {
		panic(err)
	}
	defer userRepo.Close()

	if err := userRepo.AutoMigrate(); err != nil {
		panic(">>> main > Failed at database migration! Details: " + err.Error())
	}

	staticCtrl := controllers.NewStatic()
	usersCtrl := controllers.NewUsers(userRepo)

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
