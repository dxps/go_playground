package main

import (
	"log"
	"net/http"

	"devisions.org/go-reddit/store/postgres"
	"devisions.org/go-reddit/web"
)

func main() {

	dsn := "postgres://go-reddit:secret@localhost:54326/go-reddit?sslmode=disable"

	store, err := postgres.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	csrfKey := []byte("01234567890123456789012345678901") // a 32-bytes long CSRF key

	h := web.NewHandler(store, sessions, csrfKey)
	log.Println(">>> Starting the Web server ...")
	log.Fatal(http.ListenAndServe(":3000", h))

}
