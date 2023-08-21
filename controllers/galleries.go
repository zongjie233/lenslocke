package controllers

import (
	"github.com/zongjie233/lenslocked/models"
	"net/http"
)

type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

// New returns a new Gallery entity and
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}
