package controllers

import (
	"github.com/gitalek/gogal/models"
	"github.com/gitalek/gogal/views"
)

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

func NewGalleries(gs models.GalleryService) (*Galleries, error) {
	viewNew, err := views.NewView("bootstrap", "galleries/new")
	if err != nil {
		return nil, err
	}
	return &Galleries{
		New: viewNew,
		gs: gs,
	}, nil
}
