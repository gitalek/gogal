package main

import (
	"fmt"
	"github.com/gitalek/gogal/controllers"
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

	staticC, err := controllers.NewStatic()
	must(err)
	usersC, err := controllers.NewUsers(services.User)
	must(err)
	r := mux.NewRouter()
	r.Handle("/", staticC.Home)
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(err404)
	log.Fatal(http.ListenAndServe(":3000", r))
}
