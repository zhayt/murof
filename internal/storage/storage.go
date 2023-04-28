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
	CreatePost(ctx context.Context, post model.Post) error
	GetAllPosts(ctx context.Context) ([]model.Post, error)
	GetPostByID(ctx context.Context, postID int) (*model.Post, error)
	GetPostsByCategory(ctx context.Context, categories []string) ([]model.Post, error)
	GetUserPosts(ctx context.Context, userID int) ([]model.Post, error)
	GetUserCommentPosts(ctx context.Context, userID int) ([]model.Post, error)
	GetUserLikedPosts(ctx context.Context, userID int) ([]model.Post, error)
	UpdatePost(ctx context.Context, post model.Post, postID int) error
	DeletePost(ctx context.Context, postID int) error
}

type ICommentStorage interface {
	CreateComment(ctx context.Context, comment model.Comment) error
	GetCommentByPostID(ctx context.Context, postID int) (*[]model.Comment, error)
}

type Storage struct {
}
