package components

import (
	"fmt"
	"log/slog"
	"strconv"

	omap "github.com/elliotchance/orderedmap/v3"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type item struct {
	id    int
	value string
}

func (i item) String() string {
	return fmt.Sprintf("{id=%d, value=%s}", i.id, i.value)
}

type DndList struct {
	app.Compo

	items    *omap.OrderedMap[int, string] // the items received by the component
	dragging bool                          // just for presentation
	sIdx     int                           // drag source index
	sItem    item                          // drag source item
	tIdx     int                           // drag target index
	tItem    item                          // drag target item
}

func NewDndList(items *omap.OrderedMap[int, string]) *DndList {
	return &DndList{
		items: items,
		sIdx:  -1,
		tIdx:  -1,
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
					DataSets(map[string]any{
						"itemidx": i,
						"itemid":  itemsList[i].id,
					}).
					Draggable(true).
					Text(itemsList[i].value).
					OnDragStart(func(ctx app.Context, e app.Event) {
						d.dragging = true
						d.sIdx = atoi(ctx.JSSrc().Get("dataset").Get("itemidx").String())
						id := atoi(ctx.JSSrc().Get("dataset").Get("itemid").String())
						val, _ := d.items.Get(id)
						d.sItem = item{
							id:    id,
							value: val,
						}
						slog.Debug("OnDragStart", "sIdx", i, "sItem", d.sItem)
					}).
					OnDragOver(func(ctx app.Context, e app.Event) {
						if i != d.sIdx {
							d.tIdx = atoi(ctx.JSSrc().Get("dataset").Get("itemidx").String())
							id := atoi(ctx.JSSrc().Get("dataset").Get("itemid").String())
							val, _ := d.items.Get(id)
							d.tItem = item{
								id:    id,
								value: val,
							}
						}
					}).
					OnDragEnd(func(ctx app.Context, e app.Event) {
						d.dragging = false
						if d.sIdx != d.tIdx && d.sIdx != -1 && d.tIdx != -1 {
							slog.Debug("OnDragEnd", "sIdx", d.sIdx, "tIdx", d.tIdx, "sItem", d.sItem, "tItem", d.tItem)
							d.reorderItems()
							d.sIdx = -1
							d.tIdx = -1
							ctx.Update()
						}
					}).
					Class("bg-gray-100 rounded-md px-3 py-2 m-4 hover:cursor-grab")
			}),
		)
}

func (d *DndList) reorderItems() {
	//
	newItems := omap.NewOrderedMap[int, string]()
	asc := d.sIdx < d.tIdx
	i := 0
	for id, val := range d.items.AllFromFront() {
		if asc {
			if i == d.tIdx {
				newItems.Set(d.tItem.id, d.tItem.value)
				newItems.Set(d.sItem.id, d.sItem.value)
			} else if i != d.sIdx {
				newItems.Set(id, val)
			}
		} else {
			if i == d.tIdx {
				newItems.Set(d.sItem.id, d.sItem.value)
				newItems.Set(d.tItem.id, d.tItem.value)
			} else if i != d.sIdx {
				newItems.Set(id, val)
			}
		}
		i++
	}
	slog.Debug("reorderItems", "newItems", toString(newItems))
	//
	d.items = newItems
}

func atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		slog.Error("atoi", "error", err)
	}
	return v
}

func toString(items *omap.OrderedMap[int, string]) string {
	s := "[ "
	for id, val := range items.AllFromFront() {
		s += fmt.Sprintf("{id=%d, val=%s} ", id, val)
	}
	s += "]"
	return s
}
