package main

import (
	"log"
	"os"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/http/api"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
)

// The default file path used for persisting the 'memstore' to disk,
// if no custom value is provided as DATA_FILE_PATH environment variable.
var DATA_FILE_PATH = "./secrets"

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
