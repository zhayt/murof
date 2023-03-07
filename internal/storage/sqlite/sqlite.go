package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"os"

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

func (s *Storage) Session(ctx context.Context, userId int) (*model.Session, error) {
	qr := `SELECT * FROM session WHERE user_id = ?`

	var session model.Session

	err := s.db.QueryRowContext(ctx, qr, userId).Scan(&session.ID, &session.UserId, &session.Token, &session.ExpirationTime)
	if err == sql.ErrNoRows {
		return nil, err
	}
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
