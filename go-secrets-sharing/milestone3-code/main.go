package main

import (
	"log"
	"os"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/http/api"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
)

// The default file path used for persisting the in-memory data to disk.
// If no custom value is provided as DATA_FILE_PATH environment variable.
var DATA_FILE_PATH = "./secrets"

func main() {

	var dataFilePath = DATA_FILE_PATH
	if val, exists := os.LookupEnv("DATA_FILE_PATH"); exists {
		dataFilePath = val
	}

	var storePass, storeSalt string
	if val, exists := os.LookupEnv("STORE_PASS"); exists {
		storePass = val
	} else {
		log.Fatalf("No STORE_PASS env var is provided.")
	}

	if val, exists := os.LookupEnv("STORE_SALT"); exists {
		storeSalt = val
	} else {
		log.Fatalf("No STORE_SALT env var is provided.")
	}

	repo, err := repo.NewRepo(dataFilePath, storePass, storeSalt)
	if err != nil {
		log.Fatalf("Repo init error: %v", err)
	}

	secrets := domain.NewSecrets(repo)
	httpApi := api.NewHttpApi(secrets)

	httpApi.Start()
}
