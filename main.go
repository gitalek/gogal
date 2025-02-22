package main

import (
	"fmt"
	"github.com/gitalek/gogal/controllers"
	"github.com/gitalek/gogal/middleware"
	"github.com/gitalek/gogal/models"
	"github.com/gitalek/gogal/rand"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	cfg := DefaultConfig()
	dbCfg := DefaultPostgresConfig()
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()

	// Controllers section.
	staticC, err := controllers.NewStatic()
	must(err)
	usersC, err := controllers.NewUsers(services.User, r)
	must(err)
	galleriesC, err := controllers.NewGalleries(services.Gallery, services.Image, r)
	must(err)

	// Middlewares section.
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)
	editGallery := requireUserMw.ApplyFn(galleriesC.Edit)
	updateGallery := requireUserMw.ApplyFn(galleriesC.Update)
	deleteGallery := requireUserMw.ApplyFn(galleriesC.Delete)
	indexGallery := requireUserMw.ApplyFn(galleriesC.Index)
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
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET").
		Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")
	r.HandleFunc("/galleries", indexGallery).Methods("GET").
		Name(controllers.IndexGalleries)
	r.HandleFunc(
		"/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete),
	).Methods("POST")
	// Image routes
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).
		Methods("POST")
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", imageHandler))
	// Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	csrfKey, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(csrfKey, csrf.Secure(cfg.IsProd()))
	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r))))
}
