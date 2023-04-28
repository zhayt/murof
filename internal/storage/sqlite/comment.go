package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
)

type CommentStorage struct {
	db *sql.DB
}

func (r *CommentStorage) CreateComment(ctx context.Context, comment model.Comment) error {
	qr := `INSERT INTO comment (postId, userId, author, text, date) VALUES ($2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, qr, comment.PostId, comment.UserId, comment.Author, comment.Text, comment.Date)
	if err != nil {
		return fmt.Errorf("couldn't create comment: %w", err)
	}
	return nil
}

func (r *CommentStorage) GetCommentByPostID(ctx context.Context, postID int) (*[]model.Comment, error) {
	qr := `SELECT * FROM comment WHERE postId=$1`

	rows, err := r.db.QueryContext(ctx, qr, postID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get comment by post id: %w", err)
	}

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Author, &comment.Text, &comment.Like, &comment.Dislike, &comment.Date); err != nil {
			return nil, fmt.Errorf("couldn't get comment by post id, scan error: %w", err)
		}

		comments = append(comments, comment)
	}

	return &comments, nil
}
