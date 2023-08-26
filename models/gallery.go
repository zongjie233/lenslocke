package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	Path      string
	GalleryID int
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB        *sql.DB
	ImagesDir string
}

// Create a new Gallery with the given title and userID
func (gs *GalleryService) Create(title string, userID int) (*Gallery, error) {
	// Create a new Gallery
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}
	// Insert the Gallery into the database
	row := gs.DB.QueryRow(`
			INSERT 
			INTO
				galleries
				(title,user_id)    
			VALUES
				($1,$2) RETURNING id;`, gallery.Title, gallery.UserID)
	// Scan the Gallery ID from the database
	err := row.Scan(&gallery.ID)
	// If an error occurs, return an error
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	// Return the Gallery
	return &gallery, nil
}

// do some changes
// Get a Gallery by ID
func (gs *GalleryService) ByID(id int) (*Gallery, error) {
	// Create a new Gallery
	gallery := Gallery{
		ID: id,
	}
	// Query the Gallery from the database
	row := gs.DB.QueryRow(`
			SELECT
				title,
				user_id
			FROM
				galleries
			WHERE
				id = $1;`, gallery.ID)
	// Scan the Gallery from the database
	err := row.Scan(&gallery.Title, &gallery.UserID)
	// If an error occurs, return an error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get gallery by id: %w", err)
	}
	// Return the Gallery
	return &gallery, nil
}

// ByUserID Get a Gallery by UserID 为了避免不必要的复制和更好地处理数据，返回指针类型的切片是更好的选择。
func (gs *GalleryService) ByUserID(userID int) ([]*Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT
			id,
			title 
		FROM
			galleries 
		WHERE
			user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("get galleries by user id: %w", err)
	}

	var galleries []*Gallery
	for rows.Next() {
		gallery := Gallery{UserID: userID}
		if err = rows.Scan(&gallery.ID, &gallery.Title); err != nil {
			return nil, fmt.Errorf("get galleries by user id: %w", err)
		}
		galleries = append(galleries, &gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("get galleries by user id: %w", err)
	}
	return galleries, nil
}

// Update a Gallery
func (gs *GalleryService) Update(gallery *Gallery) error {
	// Update the Gallery in the database
	_, err := gs.DB.Exec(`
		UPDATE
			galleries
		SET
			title = $2
		WHERE
			id = $1;`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

// Delete a Gallery
func (gs *GalleryService) Delete(id int) error {
	// Delete the Gallery in the database
	_, err := gs.DB.Exec(`
		DELETE FROM
			galleries
		WHERE
			id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(gs.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images : %w", err)
	}

	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, gs.extensions()) {
			images = append(images, Image{
				Path:      file,
				Filename:  filepath.Base(file),
				GalleryID: galleryID,
			})
		}
	}
	return images, nil
}

// search one image
func (gs *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(gs.galleryDir(galleryID), filename)
	// 检查文件是否存在
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image : %w", err)
	}
	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (gs *GalleryService) DeleteImage(galleryID int, filename string) error {
	// Get the images from the gallery service
	image, err := gs.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

// search all images

func (gs *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func hasExtension(file string, extensions []string) bool {
	// Loop through each extension in the array
	for _, ext := range extensions {
		// Convert the file string to lowercase
		file = strings.ToLower(file)
		// Convert the extension string to lowercase
		ext = strings.ToLower(ext)
		// Check if the file extension matches the extension
		if filepath.Ext(file) == ext {
			// If it does, return true
			return true
		}
	}
	// If the file extension does not match any of the extensions, return false
	return false
}

// galleryDir
func (gs *GalleryService) galleryDir(id int) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}
