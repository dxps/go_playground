package server

import (
	"go-app_files-mgmt/internal/shared/http/client"
	"go-app_files-mgmt/internal/ui/pages"
	"go-app_files-mgmt/internal/ui/uiroutes"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// initRoutes registers the UI routes.
func initRoutes(apiClient *client.ApiClient) {

	app.Route(uiroutes.Home, func() app.Composer { return &pages.HomePage{} })
	app.Route(uiroutes.Files, func() app.Composer { return pages.NewFilesPage(apiClient) })

}
