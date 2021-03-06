package controllers

import "github.com/gitalek/gogal/views"

func NewStatic() (*Static, error) {
	homeView, err := views.NewView("bootstrap", "static/home")
	if err != nil {
		return nil, err
	}
	contactView, err := views.NewView("bootstrap", "static/contact")
	if err != nil {
		return nil, err
	}
	return &Static{
		Home:    homeView,
		Contact: contactView,
	}, nil
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
