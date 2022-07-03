package api

import (
	"log"
	"net/http"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
)

const httpAddress = ":9001"

type HttpApi struct {
	server  *http.ServeMux
	secrets *domain.Secrets
}

func NewHttpApi(secrets *domain.Secrets) HttpApi {

	a := HttpApi{
		server:  http.NewServeMux(),
		secrets: secrets,
	}
	a.server.HandleFunc("/secrets/", a.getSecretHandler)
	a.server.HandleFunc("/secrets", a.addSecretHandler)
	a.server.HandleFunc("/healthcheck", a.healthcheckHandler)
	return a
}

func (a *HttpApi) Start() {

	log.Printf("Starting listening on %v ...\n", httpAddress)
	log.Fatalf("Error in http.ListenAndServe: %v", http.ListenAndServe(httpAddress, a.server))
}
