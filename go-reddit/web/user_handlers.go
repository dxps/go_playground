package web

import (
	"html/template"
	"net/http"

	goreddit "devisions.org/go-reddit"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *UserHandlers) ShowRegister() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/user_register.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, struct {
			SessionData
			CSRF template.HTML
		}{GetSessionData(h.sessions, r.Context()), csrf.TemplateField(r)})
	}
}

func (h *UserHandlers) RegisterSubmit() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		form := UserRegistrationForm{
			Username:      r.FormValue("username"),
			Password:      r.FormValue("password"),
			UsernameTaken: false,
		}
		if _, err := h.store.GetUserByUsername(form.Username); err == nil {
			form.UsernameTaken = true
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}
		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.store.SaveUser(&goreddit.User{
			ID:       uuid.New(),
			Username: form.Username,
			Password: string(password),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessions.Put(r.Context(), "flash", "Your registration was successful. Please log in.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
