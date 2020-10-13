package view

import (
	"fmt"
	"net/http"

	"devisions.org/andon-go/webapp/config"
)

// RegisterStaticHandlers registers HTTP handlers that will serve static
// content such as CSS and JavaScript files
func RegisterStaticHandlers() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(config.Get().StaticRoot))))
	fmt.Println(">>> Registered handler for '/static/'")
}
