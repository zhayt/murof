package handler

import (
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"html/template"
	"log"
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
	user := r.Context().Value(model.CtxUserKey).(model.User)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/" {
			errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []model.Post
		posts, err := h.service.Post.ShowAllPosts()
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/homepage.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		var category []string
		if r.FormValue("category"+string('5')) != "" {
			var posts []model.Post
			posts, err := h.service.Post.ShowAllPosts()
			if err != nil {
				log.Fatal(err)
			}
			display := Display{
				Username: user.Username,
				Posts:    posts,
			}
			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				log.Fatal(err)
			}
			temp.Execute(w, display)
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
				log.Fatal(err)
			}
			display := Display{
				Username: user.Username,
				Posts:    posts,
				Category: category,
			}
			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				log.Fatal(err)
			}
			temp.Execute(w, display)
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
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
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
		if r.URL.Path != "/myPosts" {
			errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []model.Post
		posts, err := h.service.Post.ShowMyPosts(user.Id)
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/myposts.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))
			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
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
		if r.URL.Path != "/myCommentPosts" {
			errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []model.Post
		posts, err := h.service.Post.ShowMyCommentPosts(user.Id)
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/mycommentposts.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))
			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
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
		if r.URL.Path != "/myLikedPosts" {
			errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []model.Post
		posts, err := h.service.Post.ShowMyLikedPosts(user.Id)
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/mylikedposts.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))
			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	}
}
