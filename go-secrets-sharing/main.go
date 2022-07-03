package main

import (
	"os"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/http/api"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
	log "github.com/sirupsen/logrus"
)

// The default file used for persisting the 'memstore' to disk,
// if no custom value has been defined as an environment variable.
var DATA_FILE_PATH = "./secrets.data"

func main() {

	var dataFilePath = DATA_FILE_PATH
	if val, exists := os.LookupEnv("DATA_FILE_PATH"); exists {
		dataFilePath = val
	}

	repo, err := repo.NewRepo(dataFilePath)
	if err != nil {
		log.Fatalf("Repo init error: %v", err)
	}

	secrets := domain.NewSecrets(repo)
	httpApi := api.NewHttpApi(secrets)

	httpApi.Start()
}
