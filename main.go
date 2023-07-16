package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>welcome to my big sssssss site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Contact page</h1>")
}

func main() {
	// var router Router
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not foune", http.StatusNotFound)
	})
	fmt.Println("starting the server on :3000....")
	http.ListenAndServe(":3000", r)
}
