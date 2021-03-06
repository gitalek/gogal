package controllers

import (
	"fmt"
	"github.com/gitalek/gogal/views"
	"github.com/gorilla/schema"
	"net/http"
)

type Users struct {
	NewView *views.View
}

func NewUsers() (*Users, error) {
	v, err := views.NewView("bootstrap", "views/users/new.gohtml")
	if err != nil {
		return nil, err
	}
	return &Users{NewView: v}, nil
}

// New processes the GET /new route
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create processes the POST /signup route
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()
	var form SignupForm
	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
