package usercontroller

import (
	"context"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"log"
	"net/http"
)

type userProvider interface {
	CreateUser(ctx context.Context, user *model.User) error
	User(ctx context.Context, email, password string) (*model.User, error)
}

type UserController struct {
	userProvider userProvider
}

func NewUserController(provider userProvider) *UserController {
	return &UserController{
		userProvider: provider,
	}
}

func (c *UserController) SignUpForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is sign-up form!"))
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Write([]byte("can't parse form"))
		return
	}

	email, ok := r.Form["email"]
	if !ok {
		log.Printf("Sign Up: Parse Form: email field not found")
		return
	}

	username, ok := r.Form["username"]
	if !ok {
		log.Printf("Sign Up: Parse Form: username field not found")
		return
	}

	password, ok := r.Form["password"]
	if !ok {
		log.Printf("Sign Up: Parse Form: password field not found")
		return
	}

	user := &model.User{
		Name:         username[0],
		Email:        email[0],
		PasswordHash: password[0],
	}

	if err := c.userProvider.CreateUser(context.Background(), user); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func (c *UserController) LoginForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is sign-in form!"))
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Write([]byte("can't parse form"))
		return
	}

	email, ok := r.Form["email"]
	if !ok {
		log.Printf("Sign Up: Parse Form: email field not found")
		return
	}

	password, ok := r.Form["password"]
	if !ok {
		log.Printf("Sign Up: Parse Form: password field not found")
		return
	}

	user, err := c.userProvider.User(context.Background(), email[0], password[0])
	if err != nil {
		log.Println(err)
		w.Write([]byte("server or client error"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.Token,
		Expires:  user.ExpirationTime,
		Path:     "/",
		HttpOnly: true,
	})

	w.Write([]byte("You are sign in successfully"))
}
