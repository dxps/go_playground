//go:build js && wasm

package ui

import (
	"fmt"
	"go-app_files-mgmt/internal/ui/infra"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func InitPWAClientSide(apiPort int) {

	apiClient := infra.NewApiClient(fmt.Sprintf("http://localhost:%d", apiPort))
	initAppRoutesClientSide(apiClient)
	app.RunWhenOnBrowser()
}
