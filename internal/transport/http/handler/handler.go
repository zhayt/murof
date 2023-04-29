package handler

import (
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
	l       *logger.Logger
}

func NewHandler(s *service.Service, l *logger.Logger) *Handler {
	return &Handler{
		service: s,
		l:       l,
	}
}

type Display struct {
	Username string
	Posts    []model.Post
	Category []string
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	user := r.Context().Value(model.CtxUserKey).(model.User)

	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/homepage.html")
		if err != nil {
			h.l.Error.Println("Parse file error:")
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var posts []model.Post
		posts, err = h.service.Post.ShowAllPosts()
		if err != nil {
			h.l.Error.Printf("Show all posts error: %s", err.Error())
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.l.Info.Println("All post founded, count %v", len(posts))

		display := Display{
			Username: user.Username,
			Posts:    posts,
		}

		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var category []string

		if r.FormValue("category"+string('5')) != "" {
			var posts []model.Post

			posts, err := h.service.Post.ShowAllPosts()
			if err != nil {
				h.l.Error.Printf("show all post error: %s", err.Error())
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			display := Display{
				Username: user.Username,
				Posts:    posts,
			}

			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if err = temp.Execute(w, display); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			return
		}
		for i := '0'; i <= '4'; i++ {
			if r.FormValue("category"+string(i)) != "" {
				category = append(category, r.FormValue("category"+string(i)))
			}
		}
		if len(category) != 0 {
			// fmt.Println(category)
			posts, err := h.service.Post.GetPostsByCategoty(category)
			if err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			display := Display{
				Username: user.Username,
				Posts:    posts,
				Category: category,
			}

			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if err = temp.Execute(w, display); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))
			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Like.SetPostLike(like); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}

	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) MyPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)

	var empty model.User

	if user == empty {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

	switch r.Method {
	case http.MethodGet:

		var posts []model.Post
		posts, err := h.service.Post.ShowMyPosts(user.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/myposts.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))
			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) MyCommentPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)
	var empty model.User
	if user == empty {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		var posts []model.Post

		posts, err := h.service.Post.ShowMyCommentPosts(user.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		display := Display{
			Username: user.Username,
			Posts:    posts,
		}

		temp, err := template.ParseFiles("ui/mycommentposts.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))

			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Like.SetPostLike(like); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))

			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) MyLikedPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)
	var empty model.User
	if user == empty {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

	switch r.Method {
	case http.MethodGet:
		var posts []model.Post

		posts, err := h.service.Post.ShowMyLikedPosts(user.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		display := Display{
			Username: user.Username,
			Posts:    posts,
		}

		temp, err := template.ParseFiles("ui/mylikedposts.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))

			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Like.SetPostLike(like); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))

			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
