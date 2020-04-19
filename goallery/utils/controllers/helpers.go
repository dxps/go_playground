package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
)

// ParseForm uses gorilla's schema package to
// decode a form submit request into a structure.
// dst should be a pointer to the destination structure
// that is populated by this parsing mechanism.
func ParseForm(r *http.Request, dst interface{}) error {

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
