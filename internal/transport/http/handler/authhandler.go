package handler

import (
	"errors"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"html/template"
	"net/http"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/signin.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, temp); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("psw")
		user, err := h.service.GenerateToken(email, password)
		if err != nil {
			if errors.Is(err, service.InvalidDate) {
				errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: user.Token,
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/signup.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, temp); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var signUp model.User
		email := r.FormValue("email")
		password := r.FormValue("psw")
		username := r.FormValue("username")
		repeatpass := r.FormValue("repeatspw")

		if repeatpass != password {
			errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else {
			signUp = model.User{
				Username: username,
				Login:    email,
				Password: password,
			}

		}

		if err := h.service.CreateUser(signUp); err != nil {
			if errors.Is(err, service.InvalidDate) {
				errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: "",
		Path:  "/",
	})
	/*user := r.Context().Value(ctxUserKey).(models.User)
	if err := h.service.DeleteToken(user.Token); err != nil {
		log.Fatal("delete token in logout", err)
	}*/
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
