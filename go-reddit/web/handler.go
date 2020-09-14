package web

import (
	"html/template"
	"net/http"

	goreddit "devisions.org/go-reddit"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	*chi.Mux
	store goreddit.Store
}

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)

	h.Get("/", h.HomeHandler())
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.GetThreadsHandler())
		r.Get("/new", h.NewThreadHandler())
		r.Post("/", h.SaveThreadHandler())
		r.Get("/{id}", h.ShowThreadHandler())
		r.Post("/{id}/delete", h.DeleteThreadHandler())
		r.Get("/{id}/new", h.ShowCreatePostHandler())
		r.Post("/{id}", h.SavePostHandler())
		r.Get("/{threadID}/{postID}", h.ShowPostsHandler())
	})

	return h
}

func (h *Handler) HomeHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, nil)
	}
}

func (h *Handler) GetThreadsHandler() http.HandlerFunc {
	// init
	type data struct {
		Threads []goreddit.Thread
	}
	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html", "web/templates/threads.html"))
	// handler logic
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := h.store.GetThreads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, data{Threads: ts})
	}
}

func (h *Handler) ShowThreadHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/thread.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, nil)
	}
}

func (h *Handler) NewThreadHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, nil)
	}
}

func (h *Handler) SaveThreadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.SaveThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *Handler) DeleteThreadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *Handler) ShowCreatePostHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, nil)
	}
}

func (h *Handler) ShowPostsHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/posts.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, nil)
	}
}

func (h *Handler) SavePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.SaveThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
