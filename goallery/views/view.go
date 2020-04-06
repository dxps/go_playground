package views

import "html/template"

// NewView creates a view based on the provided templates.
// It also appends the common template files.
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/bootstrap.gohtml",
		"views/layouts/navbar.gohtml",
		"views/layouts/footer.gohtml",
	)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View represents is returned as response to be presented in the browser.
type View struct {
	Template *template.Template
	Layout   string
}
