package web

import (
	"fmt"
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
		r.Post("/{threadID}/{postID}", h.SaveCommentHandler())
	})
	h.Get("/comments/{id}/vote", h.VoteCommentHandler())

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
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		t, err := h.store.GetThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ps, err := h.store.GetPostsByThread(t.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		data := struct {
			Thread goreddit.Thread
			Posts  []goreddit.Post
		}{t, ps}
		_ = tmpl.Execute(w, data)
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
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := h.store.GetThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, struct{ Thread goreddit.Thread }{t})
	}
}

func (h *Handler) ShowPostsHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		postIDStr := chi.URLParam(r, "postID")
		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := h.store.GetThread(threadID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, err := h.store.GetPost(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cs, err := h.store.GetCommentsByPost(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, struct {
			Thread   goreddit.Thread
			Post     goreddit.Post
			Comments []goreddit.Comment
		}{t, p, cs})
	}
}

func (h *Handler) SavePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		content := r.FormValue("content")

		idStr := chi.URLParam(r, "id")
		tid, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := goreddit.Post{
			ID:       uuid.New(),
			ThreadID: tid,
			Title:    title,
			Content:  content,
		}
		if err := h.store.SavePost(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/threads/%s/%s", tid, p.ID), http.StatusFound)
	}
}

func (h *Handler) SaveCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := r.FormValue("content")
		pidStr := chi.URLParam(r, "postID")
		pid, err := uuid.Parse(pidStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.store.SaveComment(&goreddit.Comment{
			ID:      uuid.New(),
			PostID:  pid,
			Content: content,
			Votes:   0,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (h *Handler) VoteCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c, err := h.store.GetComment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			c.Votes++
		} else if dir == "down" {
			c.Votes--
		}

		if err := h.store.UpdateComment(&c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
