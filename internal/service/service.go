package service

import (
	"context"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service/userservice"
)

type storage interface {
	CreateUser(ctx context.Context, user *model.User) error
	User(ctx context.Context, email string) (*model.User, error)
	CrateSession(ctx context.Context, session *model.User) error
	Session(ctx context.Context, userId int) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
}

type Service struct {
	UserService *userservice.UserService
}

func NewUserService(storage storage) *Service {
	return &Service{
		UserService: userservice.NewUserService(storage),
	}
}
