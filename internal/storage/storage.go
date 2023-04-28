package storage

import (
	"context"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
)

type Storage interface {
	CreateUser(ctx context.Context, user *model.User) error
	User(ctx context.Context, email string) (*model.User, error)
	CrateSession(ctx context.Context, session *model.Session) error
	Session(ctx context.Context, token string) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
}
