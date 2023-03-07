package controller

import (
	"github.com/zhayt/clean-arch-tmp-forum/internal/controller/usercontroller"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"net/http"
)

type Controller struct {
	UserController *usercontroller.UserController
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		UserController: usercontroller.NewUserController(service.UserService),
	}
}

func (c *Controller) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}
