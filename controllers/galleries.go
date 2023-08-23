package controllers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5"

	"github.com/zongjie233/lenslocked/context"
	"github.com/zongjie233/lenslocked/models"
	"math/rand"
	"net/http"
	"strconv"
)

type Galleries struct {
	Templates struct {
		New   Template
		Show  Template
		Edit  Template
		Index Template
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

// show galleries by id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var data struct {
		ID     int
		Title  string
		Images []string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		catImageURL := fmt.Sprintf("https://placekitten.com/%d/%d", w, h)
		data.Images = append(data.Images, catImageURL)
	}
	g.Templates.Show.Execute(w, r, data)
}

// edit galleries
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
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
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
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

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	// Create a type to store the gallery information
	type Gallery struct {
		ID    int
		Title string
	}

	// Create a struct to store the gallery information
	var data struct {
		Galleries []Gallery
	}
	// Get the user from the request context
	user := context.User(r.Context())
	// Get the gallery information from the gallery service
	galleries, err := g.GalleryService.ByUserID(user.ID)
	// If there is an error, return an internal server error
	if err != nil {
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	// Loop through the gallery information and append it to the data struct
	for _, g := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    g.ID,
			Title: g.Title,
		})
	}

	// Execute the index template with the data struct
	g.Templates.Index.Execute(w, r, data)
}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

type galleryOpt func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

// galleryByID
func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	// Convert the URL parameter "id" to an integer
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		// If the ID is not an integer, return a 404 error
		http.Error(w, "invalid ID", http.StatusNotFound)
		return nil, err
	}
	// Get the gallery by the given ID
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		// If the gallery is not found, return a 404 error
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "gallery not found", http.StatusNotFound)
			return nil, err
		}
		// If something else goes wrong, return an internal server error
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return nil, err
	}
	// Iterate through the gallery options and call the function with the gallery and any errors
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	// Return the gallery and any errors
	return gallery, nil
}

// userMustOwnGallery
func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "you are not allowed to this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}
	return nil
}
