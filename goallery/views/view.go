package views

import "html/template"

// NewView creates a view based on the provided templates.
// It also appends the common template files.
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
	}
}

// View represents is returned as response to be presented in the browser.
type View struct {
	Template *template.Template
}
