package server

import (
	"fmt"
	"go-app_files-mgmt/internal/shared/http/client"
	"log"
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

// InitAndStartWebUiClientSide sets up the UI in the "client-side" (the PWA that lives in the browser).
func InitAndStartWebUiClientSide(uiPort, apiPort int) *http.Server {

	apiClient := client.NewApiClient(fmt.Sprintf("http://localhost:%d", apiPort))
	initRoutes(apiClient)

	app.RunWhenOnBrowser()

	uiSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", uiPort),
		Handler: newAppHandler(),
	}

	go func() {
		if err := uiSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return &uiSrv
}

// InitAndStartWebUiServerSide sets up the UI in the "server-side" (for server-side rendering of the UI).
func InitAndStartWebUiServerSide(uiPort, apiPort int) *http.Server {

	apiClient := client.NewApiClient(fmt.Sprintf("http://localhost:%d", apiPort))
	initRoutes(apiClient)

	app.RunWhenOnBrowser()

	uiSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", uiPort),
		Handler: newCustomHandler(),
	}

	go func() {
		if err := uiSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return &uiSrv
}
