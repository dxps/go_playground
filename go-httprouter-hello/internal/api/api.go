package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type API struct {
	httpServer *http.Server
}

func NewAPI(httpHost string, httpPort int) (*API, error) {

	a := API{}
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", httpHost, httpPort),
		Handler:      a.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	a.httpServer = httpServer

	return &a, nil
}

func (api *API) Serve() error {

	log.Printf("Listening for HTTP API requests on port %s", api.httpServer.Addr)
	return api.httpServer.ListenAndServe()
}

func (api *API) Shutdown(stopCtx context.Context) error {

	return api.httpServer.Shutdown(stopCtx)
}
