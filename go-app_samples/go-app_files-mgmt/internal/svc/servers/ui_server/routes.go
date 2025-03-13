package uiserver

import (
	"go-app_files-mgmt/internal/ui/pages"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func initAppHomeRoute() {

	app.Route("/", func() app.Composer { return &pages.HomePage{} })

}
