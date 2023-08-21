package controllers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zongjie233/lenslocked/context"
	"github.com/zongjie233/lenslocked/models"
	"net/http"
	"strconv"
)

type Galleries struct {
	Templates struct {
		New  Template
		Edit Template
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

// create a new gallery
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

// edit galleries
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusNotFound)
		return
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	// check user have this gallery
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "you are not allowed to this gallery", http.StatusForbidden)
		return
	}
	data := struct {
		ID    int
		Title string
	}{
		ID:    gallery.ID,
		Title: gallery.Title,
	}
	g.Templates.Edit.Execute(w, r, data)

}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusNotFound)
		return
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	// check user have this gallery
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "you are not allowed to this gallery", http.StatusForbidden)
		return
	}
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)

}
