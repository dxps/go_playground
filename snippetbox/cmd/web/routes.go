package main

import (
	"github.com/justinas/alice"
	"net/http"
)

//
// routes is setting up the application's http routes.
//
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// middleware added in the chain.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standardMiddleware.Then(mux)
	//return app.recoverPanic(
	//	app.logRequest(
	//		secureHeaders(mux)))

}
