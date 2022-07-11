package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *API) routes() *httprouter.Router {

	r := httprouter.New()
	r.Handle(http.MethodGet, "/", a.indexHandler)
	return r
}
