package controllers

import (
	"devisions.org/goallery/models"
	"fmt"
	"log"
	"net/http"

	"devisions.org/goallery/views"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	repo      *models.UserRepo
}

// NewUsers creates the view for "new user" use case.
func NewUsers(repo *models.UserRepo) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		repo:      repo,
	}
}

// SignupForm is submited by user signup request.
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm is submited by user login request.
type LoginForm struct {
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
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.repo.Add(&user); err != nil {
		log.Printf(">>> Error trying to add user into repo: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintln(w, "User is", user)
}

// Login is used for processing the user login request.
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.repo.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		_, _ = fmt.Fprintln(w, "Invalid email address")
	case models.ErrInvalidPwd:
		_, _ = fmt.Fprintln(w, "Invalid password provided")
	case nil:
		_, _ = fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
