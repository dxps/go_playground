package main

import (
	"devisions.org/go-playground/hex1_url-shortener/api"
	sr "devisions.org/go-playground/hex1_url-shortener/repo/redis"
	"devisions.org/go-playground/hex1_url-shortener/shortener"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	repo := chooseRepo()
	service := shortener.NewShortUrlService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		port := httpPort()
		fmt.Printf("Listening on port %s ...\n", port)
		errs <- http.ListenAndServe(port, r)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated by %s", <-errs)
}

func chooseRepo() shortener.ShortUrlRepository {
	if os.Getenv("DB_TYPE") == "redis" {
		redisURL := os.Getenv("DB_URL")
		if redisURL == "" {
			log.Fatal("Unknown Redis URL (check if DB_URL env var is correctly provided)")
		}
		repo, err := sr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	log.Fatal("Unknown database type (check if DB_TYPE env var is correctly provided)")
	return nil
}

func httpPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
