package web

import (
	"html/template"
	"net/http"

	goreddit "devisions.org/go-reddit"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
)

type Handler struct {
	*chi.Mux
	store    goreddit.Store
	sessions *scs.SessionManager
}

func NewHandler(store goreddit.Store, sessions *scs.SessionManager, csrfKey []byte) *Handler {

	h := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	threads := ThreadsHandler{store, sessions}
	posts := PostsHandler{store, sessions}
	comments := CommentsHandler{store, sessions}
	users := UserHandlers{store, sessions}

	h.Use(middleware.Logger)
	h.Use(csrf.Protect(csrfKey, csrf.Secure(false), csrf.CookieName("csrf_token"), csrf.FieldName("csrf_token")))
	h.Use(sessions.LoadAndSave)

	h.Get("/", h.HomeHandler())

	h.Route("/threads", func(r chi.Router) {

		r.Get("/", threads.List())
		r.Get("/new", threads.New())
		r.Post("/", threads.Save())
		r.Get("/{id}", threads.Show())
		r.Post("/{id}/delete", threads.Delete())

		r.Get("/{id}/new", posts.New())
		r.Post("/{id}", posts.Save())
		r.Get("/{threadID}/{postID}", posts.Show())
		r.Get("/{threadID}/{postID}/vote", posts.Vote())

		r.Post("/{threadID}/{postID}", comments.Save())
	})

	h.Get("/comments/{id}/vote", comments.Vote())

	h.Get("/register", users.ShowRegister())
	h.Post("/register", users.RegisterSubmit())

	return h
}

func (h *Handler) HomeHandler() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		ps, err := h.store.GetPosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, struct {
			SessionData
			Posts []goreddit.Post
		}{GetSessionData(h.sessions, r.Context()), ps})
	}
}
