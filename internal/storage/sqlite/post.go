package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"strings"
)

type PostStorage struct {
	db *sql.DB
}

func (r *PostStorage) CreatePost(ctx context.Context, post model.Post) error {
	qr := `INSERT INTO post (user_id, author, title, description, date,category) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, qr, post.AuthorId, post.Author, post.Title, post.Description, post.Date, post.Category)
	if err != nil {
		return fmt.Errorf("couldn't create post: %w", err)
	}

	return nil
}

func (r *PostStorage) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM post`)
	if err != nil {
		return nil, fmt.Errorf("couldn't get all posts: %w", err)
	}

	var posts []model.Post
	for rows.Next() {
		post := new(model.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (r *PostStorage) GetPostByID(ctx context.Context, postID int) (*model.Post, error) {
	rows, err := r.db.Query(`SELECT * FROM post WHERE id=?`, postID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get post by id: %w", err)
	}

	var post model.Post
	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("couldn't get post by id: scan error: %w", err)
		}
	}

	return &post, nil
}

// Qate bolu mymkin
func (r *PostStorage) GetPostsByCategory(ctx context.Context, categories []string) ([]model.Post, error) {
	fmt.Println(strings.Join(categories, ", "))
	var posts []model.Post
	var post model.Post
	for _, category := range categories {
		rows, err := r.db.QueryContext(ctx, `SELECT * FROM post where category = ?`, category)
		if err != nil {
			return nil, fmt.Errorf("couldn't get post by category: %w", err)
		}
		for rows.Next() {
			if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
				return nil, fmt.Errorf("couldn't get post by category, scan error: %w", err)
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *PostStorage) GetUserPosts(ctx context.Context, userID int) ([]model.Post, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM post WHERE user_id=?`, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get user posts: %w", err)
	}

	var posts []model.Post
	for rows.Next() {
		post := new(model.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("couldn't get user posts, scan error: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (r *PostStorage) GetUserCommentPosts(ctx context.Context, userID int) ([]model.Post, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM post WHERE id IN (SELECT postId FROM comment WHERE userId=?)`, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get user commented post: %w", err)
	}

	var posts []model.Post
	for rows.Next() {
		post := new(model.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("couldn't get user commented post, scan error: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (r *PostStorage) GetUserLikedPosts(ctx context.Context, userID int) ([]model.Post, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM post WHERE id IN (SELECT postId FROM like WHERE userId=? AND postId IS NOT NULL)`, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get user liked post: %w", err)
	}

	var posts []model.Post
	for rows.Next() {
		post := new(model.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("couldn't get user liked post, scan error: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (r *PostStorage) UpdatePost(ctx context.Context, post model.Post, postID int) error {
	qr := `UPDATE post SET title=$1,description=$2 where id=$3;`

	_, err := r.db.ExecContext(ctx, qr, post.Title, post.Description, postID)
	if err != nil {
		return fmt.Errorf("couldn't update post: %w", err)
	}

	return nil
}

func (r *PostStorage) DeletePost(ctx context.Context, postID int) error {
	qr := `DELETE FROM post WHERE  id = ?`

	_, err := r.db.ExecContext(ctx, qr, postID)
	if err != nil {
		return fmt.Errorf("couldn't delete post: %w", err)
	}

	return nil
}
