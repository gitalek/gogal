package controllers

import (
	"fmt"
	"github.com/gitalek/gogal/models"
	"github.com/gitalek/gogal/views"
	"net/http"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

func NewUsers(us *models.UserService) (*Users, error) {
	v, err := views.NewView("bootstrap", "users/new")
	if err != nil {
		return nil, err
	}
	return &Users{NewView: v, us: us}, nil
}

// New processes the GET /new route
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create processes the POST /signup route
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintf(
		w,
		"Name is %s\nEmail is %s\nPassword is %s\n",
		form.Name, form.Email, form.Password,
	)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
