package components

import (
	"log/slog"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HTMLAttrs struct {
	app.Compo
}

func NewHTMLAttrs() *HTMLAttrs {
	return &HTMLAttrs{}
}

func (d *HTMLAttrs) Render() app.UI {

	return app.Div().
		Body(
			app.Div().
				Attr("id", "div-1-id").
				Attr("title", "div-1-title").
				DataSets(map[string]any{
					"div1-1":   "div1-1 data",
					"div1-2":   "div1-2 data",
					"div1-id":  "div1 id",
					"div1-idm": "div1 idm",
				}).
				Text("Div 1").
				OnClick(func(ctx app.Context, e app.Event) {
					slog.Info("On Div 1",
						"id", ctx.JSSrc().Get("id").String(),
						"title", ctx.JSSrc().Get("title").String(),
						"div1-1", ctx.JSSrc().Get("dataset").Get("div1-1").String(),
						"div1-id", ctx.JSSrc().Get("dataset").Get("div1-id").String(),
						"div1-idm", ctx.JSSrc().Get("dataset").Get("div1-idm").String())
				}).
				Class("bg-gray-100 rounded-md px-4 py-1 my-1 hover:bg-gray-200 hover:cursor-pointer"),
			app.P().
				Attr("id", "p-1-id").
				DataSets(map[string]any{
					"p1-1": "p-1-1 data",
					"p1-2": "p-1-2 data",
				}).
				Text("Paragraph 1").
				OnClick(func(ctx app.Context, e app.Event) {
					slog.Info("On Paragraph 1",
						"id", ctx.JSSrc().Get("id").String(),
						"data p1-1", ctx.JSSrc().Get("dataset").Get("p1-1").String(),
						"data p1-2", ctx.JSSrc().Get("dataset").Get("p1-2").String())
				}).
				Class("bg-gray-100 rounded-md px-4 py-1 my-1 hover:bg-gray-200 hover:cursor-pointer"),
			app.Span().
				Attr("id", "span-1-id").
				DataSets(map[string]any{
					"span1-1": "span-1-1 data",
					"span1-2": "span-1-2 data",
				}).
				Text("Span 1").
				OnClick(func(ctx app.Context, e app.Event) {
					slog.Info("On Span 1",
						"id", ctx.JSSrc().Get("id").String(),
						"data span1-1", ctx.JSSrc().Get("dataset").Get("span1-1").String(),
						"data span1-2", ctx.JSSrc().Get("dataset").Get("span1-2").String())
				}).
				Class("bg-gray-100 rounded-md px-4 py-1 my-1 hover:bg-gray-200 hover:cursor-pointer"),
		)
}
