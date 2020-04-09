package controllers

import "devisions.org/goallery/views"

// Static aggregates the views to "almost static" pages.
type Static struct {
	HomeView    *views.View
	ContactView *views.View
}

func NewStatic() *Static {
	return &Static{
		HomeView:    views.NewView("boostrap", "static/home"),
		ContactView: views.NewView("boostrap", "static/contact"),
	}
}
