package controllers

import (
	"github.com/gitalek/gogal/context"
	"github.com/gitalek/gogal/models"
	"github.com/gitalek/gogal/views"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	ShowGallery = "show_gallery"
)

type Galleries struct {
	New      *views.View
	ShowView *views.View
	EditView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService, r *mux.Router) (*Galleries, error) {
	viewNew, err := views.NewView("bootstrap", "galleries/new")
	if err != nil {
		return nil, err
	}
	showView, err := views.NewView("bootstrap", "galleries/show")
	if err != nil {
		return nil, err
	}
	editView, err := views.NewView("bootstrap", "galleries/edit")
	if err != nil {
		return nil, err
	}
	return &Galleries{
		New:      viewNew,
		ShowView: showView,
		EditView: editView,
		gs:       gs,
		r:        r,
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
	url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// The galleryByID method has already rendered the error for us,
		// so we just need to return here.
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// The galleryByID method has already rendered the error for us,
		// so we just need to return here.
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(
			w,
			"You do not have permission to edit this gallery",
			http.StatusNotFound,
		)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}
