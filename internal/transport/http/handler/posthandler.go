package handler

import (
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("ui/createpost.html")

	user := r.Context().Value(model.CtxUserKey).(model.User)

	if user.Username == "" {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

	switch r.Method {
	case http.MethodPost:
		title := r.FormValue("title")
		category := r.FormValue("category")
		fmt.Println(category)
		content := r.FormValue("content")
		post := model.Post{
			Title:       title,
			Description: content,
			AuthorId:    user.Id,
			Author:      user.Username,
			Date:        time.Now().Format("January 2, 2006 15:04:05"),
			Category:    category,
		}
		err := h.service.Post.CreatePost(post)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp.Execute(w, temp)
}

type DispPost struct {
	Username string
	Post     model.Post
	Comments []model.Comment
	Auth     bool
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	PostID := r.URL.Query().Get("id")
	user := r.Context().Value(model.CtxUserKey).(model.User)

	_, err := strconv.Atoi(PostID)
	if err != nil {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	post, err := h.service.GetPostByID(PostID)

	if err != nil || post.Title == "" {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	DispPost := DispPost{
		Username: user.Username,
		Post:     *post,
		Comments: []model.Comment{},
	}
	DispPost.Auth = DispPost.Username == DispPost.Post.Author
	switch r.Method {
	case http.MethodGet:
		comments, err := h.service.Comment.GetCommentByPostID(post.Id)
		if err != nil {
			log.Fatal("handler", err)
		}
		DispPost.Comments = *comments
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
			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
		}

		if r.FormValue("commentLike") != "" {
			commentId, _ := strconv.Atoi(r.FormValue("commentLike"))
			like := model.Like{
				UserID:    user.Id,
				CommentId: commentId,
			}
			if err := h.service.Like.SetCommentLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
		} else if r.FormValue("commentDislike") != "" {
			commentId, _ := strconv.Atoi(r.FormValue("commentDislike"))
			dislike := model.Dislike{
				UserID:    user.Id,
				CommentId: commentId,
			}
			if err := h.service.Dislike.SetCommentDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
		}
		if user.Token == "" {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		commentText := r.FormValue("comment")
		if strings.TrimSpace(commentText) == "" {
			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
			return
		}
		comment := model.Comment{
			UserId: user.Id,
			PostId: post.Id,
			Text:   commentText,
			Author: user.Username,
			Date:   time.Now().Format("01-02-2006 15:04:05"),
		}
		if err := h.service.Comment.CreateComment(comment); err != nil {
			fmt.Println(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	temp, err := template.ParseFiles("ui/postpage.html")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(DispPost)
	temp.Execute(w, DispPost)
}

func (h *Handler) ChangePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)
	//temp, _ := template.ParseFiles("ui/changepost.html")
	postId := r.URL.Query().Get("id")
	//fmt.Println(postId)

	oldPost, err := h.service.GetPostByID(postId)

	if err != nil {
		fmt.Println("error che tam")
	}
	//fmt.Println(oldPost)
	switch r.Method {
	case http.MethodPost:
		category := r.FormValue("category")
		//fmt.Println(category)
		title := r.FormValue("title")

		description := r.FormValue("content")

		//fmt.Println(title, description, category)
		newPost := model.Post{
			Title:       title,
			Description: description,
			Category:    category,
		}
		//fmt.Println(newPost)

		if err := h.service.Post.ChangePost(newPost, *oldPost, user); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post/?id="+postId, http.StatusSeeOther)
	}
	//temp.Execute(w, temp)
	temp, err := template.ParseFiles("ui/changepost.html")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(DispPost)
	temp.Execute(w, oldPost)

}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)
	postId := r.URL.Query().Get("id")
	post, err := h.service.GetPostByID(postId)

	if err != nil {
		fmt.Println("error che tam")
	}
	if err := h.service.Post.DeletePost(user.Id, *post); err != nil {
		fmt.Println(4)
		fmt.Println(err)
		errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
