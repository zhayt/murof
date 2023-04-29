package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"net/http"
	"time"
)

type Middleware struct {
	user service.Authorization
}

func NewMiddleware(user service.Authorization) *Middleware {
	return &Middleware{user: user}
}

func (m *Middleware) MiddleWare(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		t, err := r.Cookie("session_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.CtxUserKey, model.User{})))
			case errors.Is(err, t.Valid()):
				fmt.Println(w, http.StatusBadRequest, "invalid cookie value")
			}
			fmt.Println(w, http.StatusBadRequest, "failed to get cookie")
			return
		}
		user, err = m.user.GetUserByToken(t.Value)
		if err != nil {
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.CtxUserKey, model.User{})))
			return
		}
		if user.TokenDuration.Before(time.Now()) {
			if err := m.user.DeleteToken(user.Token); err != nil {
				fmt.Println(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.CtxUserKey, user)))
	}
}
