package main

import (
	"../../pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {

	// in-memory cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/*.page.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		// Extract the file name from the full path.
		name := filepath.Base(page)
		// Construct the template by parsing the file.
		// Before that, functions must be registered with the template.
		tmpl, err := template.New(name).Funcs(templateFunctions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Including the base 'layout' to the template.
		tmpl, err = tmpl.ParseGlob("./ui/html/*.layout.gohtml")
		if err != nil {
			return nil, err
		}
		// Including any 'partial' parts to the template.
		tmpl, err = tmpl.ParseGlob("./ui/html/*.part.gohtml")
		if err != nil {
			return nil, err
		}
		cache[name] = tmpl
	}

	return cache, nil

}

//
// humanDate generates a human friendly format of time.
//
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

//
// templateFunctions is the catalog (map) of our custom functions.
//
var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}
