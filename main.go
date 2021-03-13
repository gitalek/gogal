package main

import (
	"fmt"
	"github.com/gitalek/gogal/controllers"
	"github.com/gitalek/gogal/middleware"
	"github.com/gitalek/gogal/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 54321
	user     = "gogal"
	password = "lalala"
	dbname   = "gogal_dev"
)

func err404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>We could not find the page you "+
		"were looking for :(</h1>"+
		"<p>Please <a href=\"mailto:support@gogal.io\">email</a> us if you keep being sent to an "+
		"invalid page.</p>")
}

// must helper panics on any error
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	services, err := models.NewServices(connStr)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()

	// Controllers section.
	staticC, err := controllers.NewStatic()
	must(err)
	usersC, err := controllers.NewUsers(services.User)
	must(err)
	galleriesC, err := controllers.NewGalleries(services.Gallery, r)
	must(err)

	// Middlewares section.
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)
	editGallery := requireUserMw.ApplyFn(galleriesC.Edit)

	// Static pages routes.
	r.Handle("/", staticC.Home)
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(err404)
	// User routes.
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	// Gallery routes.
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").
		Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
