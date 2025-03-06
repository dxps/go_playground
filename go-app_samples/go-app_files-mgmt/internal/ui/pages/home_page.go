package pages

import (
	"go-app_files-mgmt/internal/ui/comps"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HomePage struct {
	app.Compo
}

func (h *HomePage) Render() app.UI {
	return app.Div().Class(
		"flex flex-col min-h-screen bg-gray-100",
	).Body(
		&comps.Navbar{},
		app.Div().
			Class("flex flex-col min-h-screen justify-center items-center text-gray-700").
			Body(
				app.Div().
					Class("w-[86px] h-[86px]").Body(app.Raw(comps.LogoIcon)),
				app.Div().Text("Files Import/Export").Class("text-3xl text-gray-400"),
				app.Div().Text("An experiment for implementing file upload and download."),
			),
	)
}
