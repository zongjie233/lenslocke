package main

import (
	stdctx "context"
	"fmt"
	"github.com/zongjie233/lenslocked/context"
	"github.com/zongjie233/lenslocked/models"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "123@123.com",
	}
	ctx = context.WithUser(ctx, &user)

	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
}
