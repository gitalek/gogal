package controllers

import (
	"fmt"
	"github.com/gitalek/gogal/context"
	"github.com/gitalek/gogal/models"
	"github.com/gitalek/gogal/views"
	"net/http"
)

type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService) (*Galleries, error) {
	viewNew, err := views.NewView("bootstrap", "galleries/new")
	if err != nil {
		return nil, err
	}
	showView, err := views.NewView("bootstrap", "galleries/show")
	if err != nil {
		return nil, err
	}
	return &Galleries{
		New: viewNew,
		ShowView: showView,
		gs:  gs,
	}, nil
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
