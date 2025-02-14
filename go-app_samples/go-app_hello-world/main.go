package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// hello is a component that displays a simple "Hello World!".
// A component is a customizable, independent, and reusable UI element.
// It is created by embedding app.Compo into a struct.
type hello struct {
	app.Compo
}

// The Render method is where the component appearance is defined.
// Here, a "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}

func main() {

	// Associate `hello` component with a path. This tells to-app what
	// component to display for that path, on both client and server-side.
	app.Route("/", func() app.Composer { return &hello{} })

	// Launch the app when in the browser.
	app.RunWhenOnBrowser()

	// Handle "/" requests.
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "A simple Hello World! example",
	})

	log.Println("Listening on http://localhost:8000 ...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
