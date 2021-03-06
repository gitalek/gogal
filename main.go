package main

import (
	"fmt"
	"github.com/gitalek/gogal/controllers"
	"github.com/gitalek/gogal/views"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(faqView.Render(w, nil))
}

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
	var err error
	homeView, err = views.NewView("bootstrap", "views/home.gohtml")
	must(err)
	contactView, err = views.NewView("bootstrap", "views/contact.gohtml")
	must(err)
	faqView, err = views.NewView("bootstrap", "views/faq.gohtml")
	must(err)
	usersC, err := controllers.NewUsers()
	must(err)
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(err404)
	log.Fatal(http.ListenAndServe(":3000", r))
}

