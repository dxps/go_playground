package controllers

import (
	"net/http"

	"devisions.org/goallery/views"
)

type Users struct {
	NewView *views.View
}

// NewUsers creates the view for "new user" use case.
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// New is the handler for rendering the (signup) form
// where a new user account can be created.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}
