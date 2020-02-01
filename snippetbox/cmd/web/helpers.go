package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// ---------------------------------------------------
//             Error response helpers.
// ---------------------------------------------------

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

//
// render is a helper method for rendering the cached templates.
//
// Args: tn - the template name, td - the template data
func (app *application) render(w http.ResponseWriter, r *http.Request, tn string, td *templateData) {

	templ, ok := app.templateCache[tn]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", tn))
		return
	}
	// A buffer to store the template rendition.
	buff := new(bytes.Buffer)
	err := templ.Execute(buff, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, _ = buff.WriteTo(w)

}

//
// addDfaultData adds the current year to the provided templateData.
//
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td

}
