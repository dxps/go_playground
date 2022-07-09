package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *API) indexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if _, err := w.Write([]byte("Hello, World!")); err != nil {
		log.Printf("[indexHandler] Response error: %v\n", err)
	}
}
