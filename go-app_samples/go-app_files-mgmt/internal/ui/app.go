//go:build js && wasm

package ui

import (
	"fmt"
	"go-app_files-mgmt/internal/ui/infra"
	"go-app_files-mgmt/internal/ui/pages"
	"go-app_files-mgmt/internal/ui/uiroutes"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// InitAndStartWebUiClientSide sets up the UI in the "client-side" (the PWA that lives in the browser).
func InitAndStartWebUiClientSide(uiPort, apiPort int) {

	apiClient := infra.NewApiClient(fmt.Sprintf("http://localhost:%d", apiPort))
	initRoutesClientSide(apiClient)
	app.RunWhenOnBrowser()
}

// initRoutesClientSide registers the UI routes.
func initRoutesClientSide(apiClient *infra.ApiClient) {

	app.Route(uiroutes.Home, func() app.Composer { return &pages.HomePage{} })
	app.Route(uiroutes.Files, func() app.Composer { return pages.NewFilesPage(apiClient) })
}
