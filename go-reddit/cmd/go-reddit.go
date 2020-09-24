package main

import (
	"log"
	"net/http"

	"devisions.org/go-reddit/store/postgres"
	"devisions.org/go-reddit/web"
)

func main() {

	store, err := postgres.NewStore("postgres://go-reddit:secret@localhost:54326/go-reddit?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	csrfKey := []byte("01234567890123456789012345678901") // a 32-bytes long CSRF key

	h := web.NewHandler(store, csrfKey)
	log.Println(">>> Starting the Web server ...")
	log.Fatal(http.ListenAndServe(":3000", h))

}
