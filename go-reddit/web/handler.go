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

	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.GetThreadsHandler())
		r.Get("/new", h.NewThreadHandler())
		r.Post("/", h.SaveThreadHandler())
		r.Post("/{id}/delete", h.DeleteThreadHandler())
	})

	return h
}

const getThreadsTemplate = `
<html>
<head>
<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
<meta content="utf-8" http-equiv="encoding">
</head>
<body>
<h2>Threads</h2>
<dl>
{{range .Threads}}
	<dt>{{.Title}}</dt>
	<dd>{{.Description}}</dd>
	<dd>
		<form action="/threads/{{.ID}}/delete" method="POST">
			<button type="submit">Delete</button>
		</form>
	</dd>
{{end}}
</dl>
<a href="/threads/new">Create thread</a>
</body>
</html>
`

func (h *Handler) GetThreadsHandler() http.HandlerFunc {
	// init
	type data struct {
		Threads []goreddit.Thread
	}
	tmpl := template.Must(template.New("").Parse(getThreadsTemplate))
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

const createThreadTemplate = `
<h2>New Thread</h2>
<form action="/threads" method="POST">
	<table>
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" /></td>
		</tr>
		<tr>
			<td>Description</td>
			<td><input type="text" name="description" /></td>
		</tr>
	</table>
	<button type="submit">Create thread</button>
</form>
`

func (h *Handler) NewThreadHandler() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(createThreadTemplate))
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
