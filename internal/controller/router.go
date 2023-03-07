package controller

import "net/http"

func (c *Controller) InitRoute() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", c.home)
	mux.Handle("/user/sigup", c.switchController([]http.HandlerFunc{c.UserController.SignUpForm, c.UserController.SignUp}))
	mux.Handle("/user/login", c.switchController([]http.HandlerFunc{c.UserController.LoginForm, c.UserController.Login}))

	return c.logRequest(mux)
}
