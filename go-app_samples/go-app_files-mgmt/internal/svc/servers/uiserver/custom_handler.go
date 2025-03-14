package uiserver

import (
	_ "embed"
	"go-app_files-mgmt/internal/common"
	"log/slog"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func newAppHandler() *app.Handler {
	return &app.Handler{
		Name:         "Files Mgmt Experiment",
		ShortName:    "XP",
		Description:  "Files Management (upload and download) Experiment",
		Title:        "Files Mgmt Experiment",
		LoadingLabel: " ",
		Icon: app.Icon{
			Default: "/web/images/favicon.svg",
			SVG:     "/web/images/favicon.svg",
		},
		BackgroundColor: "#ffffff",
		ThemeColor:      "#ffffff",
		Styles:          []string{"/web/main.css"},
		InternalURLs:    []string{"https://login.microsoftonline.com/"},
	}
}

// customHandler is an internal wrapper around the app.Handler, that overrides
// the ServeHTTP method for returning a custom version of `app.css` file.
type customHandler struct {
	app.Handler
}

func newCustomHandler() *customHandler {
	return &customHandler{Handler: *newAppHandler()}
}

func (ch *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Serve our patched `app.css` file.
	if r.URL.Path == "/app.css" {
		serveAppCss(w)
	} else if r.URL.Path == common.FilesPath {
		redirectToHomeAndTellReturn(w, r)
	} else {
		ch.Handler.ServeHTTP(w, r)
	}
}

//go:embed app.css
var appCss string

func serveAppCss(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/css")
	if _, err := w.Write([]byte(appCss)); err != nil {
		slog.Error("Failed to serve 'app.css'.", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectToHomeAndTellReturn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?backto="+r.URL.Path, http.StatusFound)
}
