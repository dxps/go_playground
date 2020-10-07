package web

import (
	"fmt"
	"html/template"
	"net/http"

	goreddit "devisions.org/go-reddit"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

type PostsHandler struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *PostsHandler) New() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/post_new.html"))
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
		_ = tmpl.Execute(w, struct {
			SessionData
			Thread goreddit.Thread
			CSRF   template.HTML
		}{GetSessionData(h.sessions, r.Context()), t, csrf.TemplateField(r)})
	}
}

func (h *PostsHandler) Show() http.HandlerFunc {

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
			SessionData
			Thread   goreddit.Thread
			Post     goreddit.Post
			Comments []goreddit.Comment
			CSRF     template.HTML
		}{GetSessionData(h.sessions, r.Context()), t, p, cs, csrf.TemplateField(r)})
	}
}

func (h *PostsHandler) Save() http.HandlerFunc {

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
		h.sessions.Put(r.Context(), "flash", "Your post has been created.")
		http.Redirect(w, r, fmt.Sprintf("/threads/%s/%s", tid, p.ID), http.StatusFound)
	}
}

func (h *PostsHandler) Vote() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "postID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.GetPost(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			p.Votes++
		} else if dir == "down" {
			p.Votes--
		}

		if err := h.store.UpdatePost(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
