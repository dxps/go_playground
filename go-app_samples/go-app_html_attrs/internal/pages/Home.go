package pages

import (
	"github.com/dxps/go_playground/tree/master/go-app_samples/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Home struct {
	app.Compo
}

func (h *Home) Render() app.UI {

	return app.Div().
		Class("flex flex-col min-h-screen bg-gray-100").
		Body(
			app.Div().
				Class("flex flex-col min-h-screen justify-center items-center drop-shadow-2xl").
				Body(
					app.Div().
						Class("bg-white rounded-md p-6 mt-8 mb-8 w-[600px] min-h-[300px]").
						Body(
							app.P().
								Class("text-2xl text-gray-500 font-medium text-center mb-4").
								Text("Drag and Drop List"),
							app.Hr().Class("text-gray-200"),
							app.P().Class("text-gray-500 text-center my-8").
								Text(`Open the browser's dev tools console and click on any of the HTML elements 
								shown below to see some of their attributes - previously set - being read.`),
							components.NewHTMLAttrs(),
						),
				),
		)
}
