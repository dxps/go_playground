package api

import (
	"context"
	"fmt"
	apiroutes "go-app_files-mgmt/internal/shared/api/api_routes"
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

func NewApiServer(
	port int,
) *ApiServer {

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

// Start starts the API Server.
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

// initRouter initializes the API Server router.
func (s *ApiServer) initRouter() {

	cors := cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},                                       //
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, //
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"}, //
		MaxAge:         600,                                                 // Maximum value not ignored by any of the major browsers.
	})
	s.router = chi.NewRouter()

	// Middlewares setup.
	s.router.Use(cors)
	s.router.Use(middleware.Logger)

	s.initRoutes()
}

// initRoutes initializes the API Server routes.
func (s *ApiServer) initRoutes() {
	s.router.Post(apiroutes.Files, s.handleFileUpload)
}
