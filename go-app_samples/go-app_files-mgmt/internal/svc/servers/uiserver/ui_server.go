package uiserver

import (
	"fmt"
	"log"
	"net/http"
)

// InitAndStartWebUiServerSide sets up the UI in the "server-side" (for server-side rendering of the UI).
func InitAndStartWebUiServerSide(uiPort, apiPort int) *http.Server {

	initAppRoutesServerSide()
	// TODO: Any non home ("/{some...}") request must be redirected to home ("/")
	//       with a query param, so that after PWA starts, the HomePage will pick
	//       up the query param and redirect (back) to the correct page.

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
