package uiserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// InitAndStartWebUiServerSide sets up the UI in the "server-side" (for server-side rendering of the UI).
func InitAndStartWebUiServerSide(uiPort, apiPort int) *http.Server {

	initAppHomeRoute()
	// TODO: Any non home ("/{some...}") request must be redirected to home ("/")
	//       with a query param, so that after PWA starts, the HomePage will pick
	//       up the query param and redirect (back) to the correct page.

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

func newSrvAppHandler() *app.Handler {
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
