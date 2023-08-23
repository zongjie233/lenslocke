package main

import (
	"fmt"
	_ "github.com/go-mail/mail/v2"
	"github.com/zongjie233/lenslocked/models"
)

func main() {
	gs := models.GalleryService{}
	fmt.Println(gs.Images(1))
}
