package main

import (
	"log"
	"net/http"

	"devisions.org/andon-go/webapp/view"
)

func main() {
	view.RegisterStaticHandlers()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
