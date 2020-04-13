package controllers

import (
	"devisions.org/goallery/models"
	"fmt"
	"log"
	"net/http"

	"devisions.org/goallery/views"
)

type Users struct {
	NewView *views.View
	repo    *models.UserRepo
}

// NewUsers creates the view for "new user" use case.
func NewUsers(repo *models.UserRepo) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		repo:    repo,
	}
}

// SignupForm is used for processing the signup request.
type SignupForm struct {
	Name     string `schema:"name"`
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

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(">>> Error parsing user create submit request body: " + err.Error())
	}

	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}
	if err := u.repo.Add(&user); err != nil {
		log.Printf(">>> Error trying to add user into repo: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintln(w, "User is", user)
}
