package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/zongjie233/lenslocked/views"
)

func excuteTemplate(w http.ResponseWriter, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("patsing template : %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	t.Excute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	tplpath := filepath.Join("templates", "home.gohtml")
	excuteTemplate(w, tplpath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {

	tplpath := filepath.Join("templates", "contact.gohtml")
	excuteTemplate(w, tplpath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {

	excuteTemplate(w, filepath.Join("templates", "faq.gohtml"))
}

func main() {
	// var router Router
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not foune", http.StatusNotFound)
	})
	fmt.Println("starting the server on :3000....")
	http.ListenAndServe(":3000", r)
}
