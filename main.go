package main

import (
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/controller"
	"github.com/zhayt/clean-arch-tmp-forum/service"
	"github.com/zhayt/clean-arch-tmp-forum/storage/sqlite"
)

func main() {
	repo := sqlite.NewStorage()
	serv := service.NewUserService(repo)

	AllHandler := controller.NewController(serv)

	fmt.Println(AllHandler.UserController.SignIn("Madi"))
	fmt.Println(AllHandler.UserController.SignUp("Madi"))
	fmt.Println(AllHandler.UserController.SignIn("Madi"))

	postId := AllHandler.PostController.CreatePost(0, "golang")
	fmt.Println(AllHandler.PostController.GetPost(0, postId))
	postId = AllHandler.PostController.CreatePost(0, "python")
	fmt.Println(AllHandler.PostController.GetPost(0, postId))
	fmt.Println(AllHandler.PostController.GetAllPosts(0))
}
