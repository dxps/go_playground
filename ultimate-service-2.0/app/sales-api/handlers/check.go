package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Using a received based mechanics to have a struct
// and a value based semantics on the receiver.

type check struct {
	log *log.Logger
}

func (c check) readiness(w http.ResponseWriter, r *http.Request) {

	status := struct{ Status string }{Status: "OK"}
	json.NewEncoder(w).Encode(status)
	c.log.Println(status)
}
