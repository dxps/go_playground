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
				Attr("data-div1", "div-1-some-data").
				Text("Div 1").
				OnClick(func(ctx app.Context, e app.Event) {
					slog.Info("On Div 1",
						"id", ctx.JSSrc().Get("id").String(),
						"title", ctx.JSSrc().Get("title").String(),
						"data-div1", ctx.JSSrc().Get("data-div1").String())
				}).
				Class("bg-gray-100 rounded-md px-4 py-1 my-1 hover:bg-gray-200 hover:cursor-pointer"),
			app.P().
				Attr("id", "p-1-id").
				Attr("title", "p-1-title").
				Attr("data-p1", "p-1-some-data").
				Text("Paragraph 1").
				OnClick(func(ctx app.Context, e app.Event) {
					slog.Info("On Paragraph 1",
						"id", ctx.JSSrc().Get("id").String(),
						"title", ctx.JSSrc().Get("title").String(),
						"data-p1", ctx.JSSrc().Get("data-p1").String())
				}).
				Class("bg-gray-100 rounded-md px-4 py-1 my-1 hover:bg-gray-200 hover:cursor-pointer"),
			app.Span().
				Attr("id", "span-1-id").
				Attr("title", "span-1-title").
				Attr("data-span1", "span-1-some-data").
				Text("Span 1").
				OnClick(func(ctx app.Context, e app.Event) {
					slog.Info("On Span 1",
						"id", ctx.JSSrc().Get("id").String(),
						"title", ctx.JSSrc().Get("title").String(),
						"data-span1", ctx.JSSrc().Get("data-span1").String())
				}).
				Class("bg-gray-100 rounded-md px-4 py-1 my-1 hover:bg-gray-200 hover:cursor-pointer"),
		)
}
