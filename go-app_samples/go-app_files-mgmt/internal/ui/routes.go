//go:build js

package ui

import (
	"go-app_files-mgmt/internal/common"
	"go-app_files-mgmt/internal/ui/infra"
	"go-app_files-mgmt/internal/ui/pages"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func initAppRoutesClientSide(apiClient *infra.ApiClient) {

	app.Route(common.HomePath, func() app.Composer { return &pages.HomePage{} })
	app.Route(common.FilesPath, func() app.Composer { return pages.NewFilesPage(apiClient) })
}
