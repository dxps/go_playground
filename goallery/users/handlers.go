package users

import (
	"fmt"
	"log"
	"net/http"

	"devisions.org/goallery/utils/controllers"
	"devisions.org/goallery/utils/rand"

	"devisions.org/goallery/views"
)

// UserHandlers contains the handlers for all user related requests.
type UserHandlers struct {
	NewView   *views.View
	LoginView *views.View
	svc       UserService
}

// NewUserHandlers creates the view for "new user" use case.
func NewUserHandlers(svc UserService) *UserHandlers {
	return &UserHandlers{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		svc:       svc,
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
func (u *UserHandlers) New(w http.ResponseWriter, _ *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is used for processing the signup form submit request.
func (u *UserHandlers) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	if err := controllers.ParseForm(r, &form); err != nil {
		panic(">>> Error parsing user create submit request body: " + err.Error())
	}
	user := User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.svc.Create(&user); err != nil {
		log.Printf(">>> Error trying to add user into repo: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		// temporary render the error message for debugging
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// Login is used for processing the user login request.
func (u *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}
	if err := controllers.ParseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.svc.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case ErrNotFound:
			_, _ = fmt.Fprintln(w, "Invalid email address")
		case ErrInvalidPwd:
			_, _ = fmt.Fprintln(w, "Invalid password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// signIn is used for signing the given user in via the remember cookie.
func (u *UserHandlers) signIn(w http.ResponseWriter, user *User) error {

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.svc.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name: "remember", Value: user.Remember,
		HttpOnly: true, // this cookie isn't accessible for scripts
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest is just a temporary test.
func (u *UserHandlers) CookieTest(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("remember")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.svc.GetByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintln(w, user)
}
