package handlers

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

type Handlers struct {
	httpPort int
	httpSrv  *http.Server
	router   http.Handler
}

func New(httpPort int, assetsFS embed.FS) *Handlers {

	h := Handlers{httpPort: httpPort}
	h.init(httpPort, assetsFS)
	return &h
}

func (h *Handlers) init(httpPort int, assetsFS embed.FS) {

	h.initRouter(assetsFS)
	hs := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: h.router,
	}
	h.httpSrv = hs
}

func (h *Handlers) initRouter(assetsFS embed.FS) {

	mux := http.NewServeMux()

	assetsHandler := http.FileServerFS(assetsFS)
	mux.Handle("GET /assets/", assetsHandler)

	h.router = mux
}

func (h *Handlers) Start() {

	slog.Info(fmt.Sprintf("Listening on port %d ...", h.httpPort))
	log.Fatal(h.httpSrv.ListenAndServe())
}
