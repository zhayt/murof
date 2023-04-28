package storage

import (
	"context"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"time"
)

type IUserStorage interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserByToken(ctx context.Context, token string) (model.User, error)
	SaveToken(ctx context.Context, login string, token string, duration time.Time) error
	DeleteToken(ctx context.Context, token string) error
}

type IPostStorage interface {
}

type ICommentStorage interface {
}

type Storage struct {
}
