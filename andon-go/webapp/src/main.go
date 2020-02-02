package main

import (
	"devisions.org/andon-go/webapp/view"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {

	// The " -" before the action ends tells the template to consume any white space
	// that exists between the end of the curly brace and the next content.
	content := `{{/* a comment */ -}}
	this is a template!`

	t := template.New("my first template")
	t, err := t.Parse(content)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, t)
	if err != nil {
		log.Fatal(err)
	}

	view.RegisterStaticHandlers()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
