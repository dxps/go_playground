package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

// View represents is returned as response to be presented in the browser.
type View struct {
	Template *template.Template
	Layout   string
}

// NewView creates a view based on the provided templates.
// It also appends the common template files.
func NewView(layout string, files ...string) *View {

	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// ServeHTTP method is used for implementing the http.Handler interface.
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// layoutFiles collects all the files found in the LayoutDir
// that matches the TemplateExt extension.
func layoutFiles() []string {

	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}
