package controllers

import (
	"fmt"
	"github.com/gitalek/gogal/context"
	"github.com/gitalek/gogal/models"
	"github.com/gitalek/gogal/views"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	IndexGalleries  = "index_galleries"
	ShowGallery     = "show_gallery"
	EditGallery     = "edit_gallery"
	maxMultiPartMem = 1 << 20 // 1 megabyte
)

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	r         *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService, is models.ImageService, r *mux.Router) (*Galleries, error) {
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
	indexView, err := views.NewView("bootstrap", "galleries/index")
	if err != nil {
		return nil, err
	}
	return &Galleries{
		New:       viewNew,
		ShowView:  showView,
		EditView:  editView,
		IndexView: indexView,
		gs:        gs,
		is:        is,
		r:         r,
	}, nil
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
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
	g.ShowView.Render(w, r, vd)
}

// POST /galleries/:id/edit
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
			http.StatusForbidden,
		)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// The galleryByID method has already rendered the error for us,
		// so we just need to return here.
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	//todo: persist to db
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery successfully updated!",
		}
	}
	g.ShowView.Render(w, r, vd)
}

// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
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
			"You do not have permission to delete this gallery",
			http.StatusForbidden,
		)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.ShowView.Render(w, r, vd)
		return
	}
	url, err := g.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultiPartMem)
	if err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	// Iterate over uploaded files
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()

		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Images successfully uploaded!",
	}
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request)  {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit this gallery or image", http.StatusForbidden)
		return
	}

	// Get the filename from the path.
	filename := mux.Vars(r)["filename"]
	// Build the Image model
	i := models.Image{
		Filename: filename,
		GalleryID: gallery.ID,
	}
	// Try to delete the image.
	err = g.is.Delete(&i)
	if err != nil {
		// Render the edit page with any errors.
		var vd views.Data
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	// If all goes well, redirect to the edit gallery page.
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}
