package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dxps/go_playground/tree/master/go-app_samples/internal/pages"
	"github.com/dxps/go_playground/tree/master/go-app_samples/internal/uiserver"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const UI_PORT = 9099

func main() {

	initLogging()

	app.Route("/", func() app.Composer { return &pages.Home{} })

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:         "go-app DnD List",
		Description:  "A go-app based Drag-n-Drop example",
		Title:        "Drag-n-Drop List | go-app Demo",
		Styles:       []string{"/web/main.css"},
		LoadingLabel: "{progress}% loaded",
		Icon: app.Icon{
			Default: "/web/images/favicon.svg",
			SVG:     "/web/images/favicon.svg",
		},
	})

	slog.Info(fmt.Sprintf("Listening on http://localhost:%d", UI_PORT))
	handler, err := uiserver.NewCustomAppHandler()
	if err != nil {
		slog.Error("Failed to init app handler", "error", err)
		return
	}
	uiSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", UI_PORT),
		Handler: handler,
	}
	if err := uiSrv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("UI Server failure", "error", err)
		return
	}
}

func initLogging() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				// Log only the filename, without the extension.
				// s.File = path.Base(s.File)
				s.File, _ = strings.CutSuffix(path.Base(s.File), ".go")
			} else
			// Remove the timestamp.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			} else
			// Remove the logging level.
			if a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	slog.SetDefault(logger)
}
