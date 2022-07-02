package main

import (
	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/http/api"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
	log "github.com/sirupsen/logrus"
)

// The default file used for persisting the 'memstore' to disk,
// if no custom value has been defined as an environment variable.
const DATA_FILE_PATH = "./secrets.data"

func main() {

	dataFilePath := DATA_FILE_PATH

	repo, err := repo.NewRepo(dataFilePath)
	if err != nil {
		log.Fatalf("Repo init error: %v", err)
	}

	secrets := domain.NewSecrets(repo)
	httpApi := api.NewHttpApi(secrets)

	httpApi.Start()
}
