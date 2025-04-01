package components

import (
	"log/slog"
	"strconv"

	omap "github.com/elliotchance/orderedmap/v3"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type item struct {
	id    int
	value string
}

type DndList struct {
	app.Compo

	items    *omap.OrderedMap[int, string] // the items received by the component
	dragging bool                          // just for presentation
	sIdx     int                           // drag source index
	tIdx     int                           // drag target index
}

func NewDndList(items *omap.OrderedMap[int, string]) *DndList {
	return &DndList{
		items: items,
	}
}

func (d *DndList) Render() app.UI {

	// Only a classic (not a custom) map is supported by `app.Range()`.
	itemsList := make([]item, 0)
	for id, val := range d.items.AllFromFront() {
		itemsList = append(itemsList, item{id, val})
	}

	return app.Div().
		Body(
			app.Range(itemsList).Slice(func(i int) app.UI {
				return app.Div().
					Class("bg-gray-100 rounded-md px-3 py-2 m-4 hover:cursor-grab").
					Attr("id", i).
					Draggable(true).
					Text(itemsList[i].value).
					OnDragStart(func(ctx app.Context, e app.Event) {
						i := atoi(ctx.JSSrc().Get("id").String())
						d.dragging = true
						d.sIdx = i
						slog.Debug("OnDragStart", "sIdx", i)
					}).
					OnDragOver(func(ctx app.Context, e app.Event) {
						if i != d.sIdx {
							i := atoi(ctx.JSSrc().Get("id").String())
							d.tIdx = i
							slog.Debug("OnDragOver", "tIdx", i)
						}
					}).
					OnDragEnd(func(ctx app.Context, e app.Event) {
						d.dragging = false
						slog.Debug("DnD result", "sIdx", d.sIdx, "tIdx", d.tIdx)
						d.reorderItems(itemsList)
						ctx.Update()
					})
			}),
		)
}

func (d *DndList) reorderItems(itemsList []item) {

	newItems := omap.NewOrderedMap[int, string]()
	asc := d.sIdx < d.tIdx
	for i, item := range itemsList {
		if (i < d.sIdx && i < d.tIdx) || (i > d.sIdx && i > d.tIdx) || (i != d.sIdx && i != d.tIdx) {
			// Copy anything outside and within the dragged range.
			newItems.Set(item.id, item.value)
			continue
		}
		if asc {
			if i == d.sIdx {
				continue
			}
			if i == d.tIdx {

				newItems.Set(d.tIdx)
				newItems.Set(sid, sval)
				continue
			}
		} else {
			if id == sid {
				newItems.Set(tid, tval)
				newItems.Set(sid, sval)
				continue
			}
			if id == tid {
				continue
			}
		}
		i++
	}

	d.items = itemsList
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
