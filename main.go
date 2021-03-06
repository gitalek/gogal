package main

import (
	"fmt"
	"github.com/gitalek/gogal/controllers"
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
	var err error
	staticC, err := controllers.NewStatic()
	must(err)
	usersC, err := controllers.NewUsers()
	must(err)
	r := mux.NewRouter()
	r.Handle("/", staticC.Home)
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(err404)
	log.Fatal(http.ListenAndServe(":3000", r))
}
