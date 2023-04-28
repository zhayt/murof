package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(driver string, dsn string) (*Storage, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init(path string) error {
	stmt, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("can't read migaration file: %w", err)
	}

	_, err = s.db.Exec(string(stmt))
	if err != nil {
		return fmt.Errorf("can't init database: %w", err)
	}

	return nil
}

func (s *Storage) CreateUser(ctx context.Context, user *model.User) error {
	qr := `INSERT INTO user (user_name, user_email, password_hash) VALUES (?, ?, ?)`

	_, err := s.db.ExecContext(ctx, qr, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	return nil
}

func (s *Storage) User(ctx context.Context, email string) (*model.User, error) {
	qr := `SELECT * FROM user WHERE user_email = ?`

	var user model.User

	err := s.db.QueryRowContext(ctx, qr, email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("can't get user by id: %w", err)
	}

	return &user, nil
}

func (s *Storage) CrateSession(ctx context.Context, user *model.User) error {
	qr := `INSERT INTO session (user_id, token, expiration_time) VALUES (?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, qr, user.ID, user.Token, user.ExpirationTime); err != nil {
		return fmt.Errorf("can't create session: %w", err)
	}

	return nil
}

func (s *Storage) SessionByUserID(ctx context.Context, userId int) (*model.Session, error) {
	qr := `SELECT * FROM session WHERE user_id = ?`

	var session model.Session

	err := s.db.QueryRowContext(ctx, qr, userId).Scan(&session.ID, &session.UserId, &session.Token, &session.ExpirationTime)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Storage) SessionByToken(ctx context.Context, token string) (*model.Session, error) {
	qr := `SELECT * FROM session WHERE token = ?`

	var session model.Session

	err := s.db.QueryRowContext(ctx, qr, token).Scan(&session.ID, &session.UserId, &session.Token, &session.ExpirationTime)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Storage) UpdateSession(ctx context.Context, session *model.Session) error {
	qr := `UPDATE session SET token = ?, expiration_time = ? WHERE user_id = ?`

	if _, err := s.db.ExecContext(ctx, qr, session.Token, session.ExpirationTime, session.UserId); err != nil {
		return fmt.Errorf("can't update user session: %w", err)
	}

	return nil
}

func (s *Storage) CreatePost(ctx context.Context, post *model.Post) (int, error) {
	qr := `INSERT INTO posts(user_id, title, content, date_creation)`

	res, err := s.db.ExecContext(ctx, qr, post.UserId, post.Title, post.Content, time.Now())
	if err != nil {
		return 0, fmt.Errorf("can't create post: %w", err)
	}

	postId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't take created post id: %w", err)
	}

	return int(postId), nil
}

func (s *Storage) Post(ctx context.Context, postId int) (*model.Post, error) {
	qr := `SELECT * FROM post WHERE id = ?`

	var post model.Post

	if err := s.db.QueryRowContext(ctx, qr, postId).Scan(&post.ID, &post.UserId, &post.Title, &post.Content, &post.Likes,
		&post.Dislikes, &post.DataCreation); err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *Storage) Posts(ctx context.Context, limit int) ([]*model.Post, error) {
	qr := `SELECT * FROM post ORDER BY date_creation LIMIT ?`

	rows, err := s.db.QueryContext(ctx, qr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*model.Post, 0, limit)
	for rows.Next() {
		p := &model.Post{}
		if err := rows.Scan(&p.ID, &p.UserId, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.DataCreation); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (s *Storage) CategoryId(ctx context.Context, category string) (int, error) {
	qr := `SELECT id FROM category WHERE category_name = ?`

	var categoryId int

	err := s.db.QueryRowContext(ctx, qr, category).Scan(&categoryId)
	if err != nil {
		return 0, err
	}

	return categoryId, nil
}

func (s *Storage) CreateCategory(ctx context.Context, category string) (int, error) {
	qr := `INSERT INTO category(category_name) VALUES (?)`

	res, err := s.db.ExecContext(ctx, qr, category)
	if err != nil {
		return 0, err
	}

	categoryId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't take created category id: %w", err)
	}

	return int(categoryId), nil
}

func (s *Storage) CreatePostCategory(ctx context.Context, postId int, categoryId int) error {
	qr := `INSERT INTO post_category(post_id, category_id) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, qr, postId, categoryId); err != nil {
		return err
	}

	return nil
}
