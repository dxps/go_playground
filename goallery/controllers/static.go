package controllers

import "devisions.org/goallery/views"

// Static aggregates the views to "almost static" pages.
type Static struct {
	HomeView    *views.View
	ContactView *views.View
}

func NewStatic() *Static {
	return &Static{
		HomeView:    views.NewView("bootstrap", "static/home"),
		ContactView: views.NewView("bootstrap", "static/contact"),
	}
}
