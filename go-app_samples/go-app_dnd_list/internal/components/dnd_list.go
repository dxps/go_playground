package components

import (
	"maps"

	omap "github.com/elliotchance/orderedmap/v3"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type DndList struct {
	app.Compo

	items *omap.OrderedMap[string, string]
}

func NewDndList(items *omap.OrderedMap[string, string]) *DndList {
	return &DndList{
		items: items,
	}
}

func (d *DndList) Render() app.UI {

	// Only a classic (not a custom) map is supported by `app.Range()`.
	items := maps.Collect(d.items.AllFromBack())

	return app.Ul().
		Class("list-none").
		Body(
			app.Range(items).Map(func(k string) app.UI {
				return app.Li().
					Class("bg-gray-100 rounded-md px-3 py-2 m-4 hover:cursor-grab").
					Draggable(true).
					Text(items[k])
			}),
		)
}
