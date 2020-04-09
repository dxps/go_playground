package controllers

import "devisions.org/goallery/views"

// Static aggregates the views to "almost static" pages.
type Static struct {
	HomeView    *views.View
	ContactView *views.View
}

func NewStatic() *Static {
	return &Static{
		HomeView:    views.NewView("boostrap", "views/static/home.gohtml"),
		ContactView: views.NewView("boostrap", "views/static/contact.gohtml"),
	}
}
