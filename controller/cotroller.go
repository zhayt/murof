package controller

import (
	"github.com/zhayt/clean-arch-tmp-forum/controller/postcontroller"
	"github.com/zhayt/clean-arch-tmp-forum/controller/usercontroller"
	"github.com/zhayt/clean-arch-tmp-forum/service"
)

type userController interface {
	SignIn(name string) bool
	SignUp(name string) int
}

type postController interface {
	CreatePost(idUser int, text string) int
	GetPost(userId, postID int) string
	GetAllPosts(userId int) []string
}

type Controller struct {
	UserController userController
	PostController postController
}

func NewController(service *service.UserService) *Controller {
	return &Controller{
		UserController: usercontroller.NewUserController(service.Storage),
		PostController: postcontroller.NewPostController(service.Storage),
	}
}
