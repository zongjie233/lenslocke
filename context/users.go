package context

import (
	"github.com/zongjie233/lenslocked/models"
	"golang.org/x/net/context"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	//返回一个新的Context，其中包含了用户模型信息
	//将用户模型关联（联接）到新的Context中，并将这个新Context作为结果返回
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		return nil
	}
	return user
}
