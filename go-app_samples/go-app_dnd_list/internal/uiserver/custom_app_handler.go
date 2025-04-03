package uiserver

import (
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/dxps/go_playground/tree/master/go-app_samples/internal/uiserver/repos"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func newAppHandler() *app.Handler {
	return &app.Handler{
		Name:         "DnD Experiment using go-app",
		ShortName:    "XP",
		Description:  "A Drag-n-Drop implementation using go-app library",
		Title:        "DnD Xp | go-app",
		LoadingLabel: "{progress}% loaded",
		Icon: app.Icon{
			Default: "/web/images/favicon.svg",
			SVG:     "/web/images/favicon.svg",
		},
		BackgroundColor: "#ffffff",
		ThemeColor:      "#ffffff",
		Styles:          []string{"/web/main.css"},
	}
}

// customAppHandler is an internal wrapper around the app.Handler, that overrides
// the `ServeHTTP` method for returning a custom version of `app.css` file.
type customAppHandler struct {
	app.Handler
	urlPathRegex      *regexp.Regexp
	compressedAppWASM []byte
	appWASMSize       string
}

func NewCustomAppHandler() (*customAppHandler, error) {

	cwasm, wasmsz, err := repos.GetGzCompressedAppWASM()
	if err != nil {
		return nil, fmt.Errorf("Failed to get the compressed app.wasm: %w", err)
	}
	return &customAppHandler{
		Handler:           *newAppHandler(),
		urlPathRegex:      regexp.MustCompile(`^/attrs`),
		compressedAppWASM: cwasm,
		appWASMSize:       fmt.Sprint(wasmsz),
	}, nil
}

func (ch *customAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Serve our patched `app.css` file.
	if r.URL.Path == "/app.css" {
		serveAppCss(w)
	} else if r.URL.Path == "/web/app.wasm" {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			ch.serveCompressedAppWASM(w)
		} else {
			ch.Handler.ServeHTTP(w, r)
		}
	} else
	// For any initial request to a non-root (i.e. non-`/`) path, redirect to the home page.
	// That's because the routes set for the server-side contains only the HomePage, being
	// the page that doesn't require an HTTP Client. This client is `syscalls/js' based,
	// and this is not supported on the server side (since GOOS is not `js').
	if ch.urlPathRegex.MatchString(r.URL.Path) {
		redirectToHomeAndTellToReturn(w, r)
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

func (ch *customAppHandler) serveCompressedAppWASM(w http.ResponseWriter) {

	headers := w.Header()
	headers.Set("Content-Encoding", "gzip")
	headers.Set("Content-Length", ch.appWASMSize)
	headers.Set("Content-Type", "application/wasm")
	if _, err := w.Write(ch.compressedAppWASM); err != nil {
		slog.Error("Failed to serve compressed 'app.wasm'.", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectToHomeAndTellToReturn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?backto="+r.URL.RequestURI(), http.StatusFound)
}
