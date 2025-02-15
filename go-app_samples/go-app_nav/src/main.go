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

	app.Route("/about", func() app.Composer { return &About{} })
	app.Route("/", func() app.Composer { return &Homepage{} })

	app.RunWhenOnBrowser()

	appHandler := &app.Handler{
		Name:        "Hello",
		Description: "A simple go-app based navigation example",
		Title:       "go-app :: Navigation",
	}

	s := http.Server{
		Addr:    ":8000",
		Handler: appHandler,
	}

	log.Println("Listening on http://localhost:8000 ...")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
