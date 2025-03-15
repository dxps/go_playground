package apiserver

import (
	"context"
	"fmt"
	"go-app_files-mgmt/internal/common"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type ApiServer struct {
	router *chi.Mux
	server *http.Server
	Port   int
}

func NewApiServer(port int) *ApiServer {

	apiSrv := ApiServer{
		Port: port,
	}
	apiSrv.initRouter()
	apiSrv.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: apiSrv.router,
	}
	return &apiSrv
}

func (s *ApiServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

// Stop gracefully shuts down the API Server.
func (s *ApiServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// initRouter initializes the API Server's router.
func (s *ApiServer) initRouter() {

	cors := cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},                                                                                    //
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}, //
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},                                              //
		MaxAge:         600,                                                                                              // Maximum value not ignored by any of the major browsers.
	})
	s.router = chi.NewRouter()

	// Middlewares setup.
	s.router.Use(cors)
	s.router.Use(middleware.Logger)

	// Routes setup.
	s.router.Post(common.FilesPath, s.handleFileUpload)
}
