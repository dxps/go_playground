package main

import (
	"log"
	"net/http"

	"devisions.org/go-reddit/store/postgres"
	"devisions.org/go-reddit/web"
)

func main() {

	store, err := postgres.NewStore("postgres://postgres:secret@localhost:54326/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	log.Println(">>> Starting the Web server ...")
	log.Fatal(http.ListenAndServe(":3000", h))

}
