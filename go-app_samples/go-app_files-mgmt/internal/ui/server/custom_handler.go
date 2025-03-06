package server

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// customHandler is an internal wrapper around the app.Handler,
// that allows us to override the ServeHTTP method for returning
// a custom version of `app.css` file.
type customHandler struct {
	app.Handler
}

func newCustomHandler() *customHandler {
	return &customHandler{
		Handler: *newAppHandler(),
	}
}

func (ch *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Serve our patched `app.css` file.
	if r.URL.Path == "/app.css" {
		ServeAppCss(w)
	} else {
		ch.Handler.ServeHTTP(w, r)
	}
}

//go:embed app.css
var appCss string

func ServeAppCss(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/css")
	if _, err := w.Write([]byte(appCss)); err != nil {
		slog.Error("GetAppCss failed to write response body.", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
