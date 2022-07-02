package main

import (
	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/http/api"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
)

const SECRETS_FILE_PATH = "secrets.data"

func main() {

	repo := repo.NewRepo()
	secrets := domain.NewSecrets(repo)
	httpApi := api.NewHttpApi(secrets)

	httpApi.Start()
}
