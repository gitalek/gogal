package controllers

import "github.com/gitalek/gogal/views"

func NewStatic() (*Static, error) {
	homeView, err := views.NewView("bootstrap", "views/static/home.gohtml")
	if err != nil {
		return nil, err
	}
	contactView, err := views.NewView("bootsrap", "views/static/contact.gohtml")
	if err != nil {
		return nil, err
	}
	return &Static{
		HomeView:    homeView,
		ContactView: contactView,
	}, nil
}

type Static struct {
	HomeView    *views.View
	ContactView *views.View
}
