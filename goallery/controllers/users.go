package controllers

import (
	"fmt"
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

// SignupForm is used for processing the signup request.
type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// New is the handler for rendering the signup page
// where a new user account can be created.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is used for processing the signup form submit request.
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintln(w, "email: ", form.Email)
	_, _ = fmt.Fprintln(w, "password: ", form.Password)

}
