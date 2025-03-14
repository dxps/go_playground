package uiserver

import (
	"go-app_files-mgmt/internal/common"
	"go-app_files-mgmt/internal/ui/pages"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func initAppHomeRoutesServerSide() {
	app.Route(common.HomePath, func() app.Composer { return &pages.HomePage{} })
}
