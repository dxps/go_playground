package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Homepage struct {
	app.Compo
}

func (h *Homepage) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Homepage"),
		app.A().Href("/about").Text("About"),
	)
}

type About struct {
	app.Compo
}

func (h *About) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("About"),
		app.A().Href("/").Text("Back to home"),
	)
}

func main() {

	app.Route("/", func() app.Composer { return &Homepage{} })
	app.Route("/about", func() app.Composer { return &About{} })

	// Launch the app when in the browser.
	app.RunWhenOnBrowser()

	// Handle "/" requests.
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "A simple Hello World! example",
		Title:       "go-app :: Hello World!",
	})

	log.Println("Listening on http://localhost:8000 ...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
