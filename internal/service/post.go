package service

import (
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"log"
	"strconv"
)

type PostService struct {
	repo repository.Post
	l    *logger.Logger
}

func NewPostService(repo repository.Post, l *logger.Logger) *PostService {
	return &PostService{
		repo: repo,
		l:    l,
	}
}

func (p *PostService) CreatePost(post model.Post) error {
	return p.repo.CreatePost(post)
}

func (p *PostService) ShowAllPosts() ([]model.Post, error) {
	posts, err := p.repo.ShowAllPosts()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) GetPostByID(id string) (*model.Post, error) {
	_, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	post, err := p.repo.GetPostByID(id)
	if err != nil {
		fmt.Printf("service: %s", err)
		return nil, err
	}
	return post, nil
}

func (p *PostService) ChangePost(newpost, oldPost model.Post, user model.User) error {
	if user.Username != oldPost.Author {
		return fmt.Errorf("Uncorrect change post author ")
	}
	if err := p.repo.ChangePost(newpost, oldPost.Id); err != nil {
		return fmt.Errorf("CHANGE :%w", err)
	}
	return nil
}

func (p *PostService) ShowMyPosts(userId int) ([]model.Post, error) {
	posts, err := p.repo.ShowMyPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) ShowMyCommentPosts(userId int) ([]model.Post, error) {
	posts, err := p.repo.ShowMyCommentPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) ShowMyLikedPosts(userId int) ([]model.Post, error) {
	posts, err := p.repo.ShowMyLikedPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) GetPostsByCategoty(category []string) ([]model.Post, error) {
	posts, err := p.repo.GetPostsByCategoty(category)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) DeletePost(userId int, post model.Post) error {
	if userId != post.AuthorId {
		return fmt.Errorf("you can't delete this post! because you aren't the author!")
	}
	if err := p.repo.DeletePost(post.Id); err != nil {
		return err
	}
	return nil
}
