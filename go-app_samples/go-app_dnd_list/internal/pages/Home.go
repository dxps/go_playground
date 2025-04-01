package pages

import (
	"github.com/dxps/go_playground/tree/master/go-app_samples/internal/components"
	omap "github.com/elliotchance/orderedmap/v3"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Home struct {
	app.Compo
}

func (h *Home) Render() app.UI {

	items := omap.NewOrderedMap[string, string]()
	items.Set("1", "First Item")
	items.Set("2", "Second Item")
	items.Set("3", "Third Item")
	items.Set("4", "Fourth Item")
	items.Set("5", "Fifth Item")

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
								Class("text-2xl text-gray-500 font-medium text-center").
								Text("Drag and Drop List"),
							app.Hr().Class("text-gray-200 mt-4 mb-8"),
							components.NewDndList(items),
						),
				),
		)
}
