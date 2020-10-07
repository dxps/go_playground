package web

import (
	"html/template"
	"net/http"

	goreddit "devisions.org/go-reddit"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

type ThreadsHandler struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *ThreadsHandler) List() http.HandlerFunc {

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
		_ = tmpl.Execute(w, struct {
			SessionData
			Threads []goreddit.Thread
		}{GetSessionData(h.sessions, r.Context()), ts})
	}
}

func (h *ThreadsHandler) Show() http.HandlerFunc {

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

		_ = tmpl.Execute(w, struct {
			SessionData
			CSRF   template.HTML
			Thread goreddit.Thread
			Posts  []goreddit.Post
		}{GetSessionData(h.sessions, r.Context()), csrf.TemplateField(r), t, ps})
	}
}

func (h *ThreadsHandler) New() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/thread_new.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, struct {
			SessionData
			CSRF template.HTML
		}{GetSessionData(h.sessions, r.Context()), csrf.TemplateField(r)})
	}
}

func (h *ThreadsHandler) Save() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		form := CreateThreadForm{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		if err := h.store.SaveThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       form.Title,
			Description: form.Description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessions.Put(r.Context(), "flash", "Your new thread has been created.")
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *ThreadsHandler) Delete() http.HandlerFunc {

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
		h.sessions.Put(r.Context(), "flash", "The thread has been deleted.")
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
