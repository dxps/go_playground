package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
)

// parseForm uses gorilla's schema package to
// decode a form submital request into a structure.
// dst should be a pointer to the destination structure
// that is populated by this parsing mechanism.
func parseForm(r *http.Request, dst interface{}) error {

	// r.ParseForm() needs to be called first to populate r.PostForm
	// and then the decoder to have values to work with.
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil

}
