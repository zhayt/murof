package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"time"
)

type UserStorage struct {
	db *sql.DB
}

func (r *UserStorage) CreateUser(ctx context.Context, user model.User) error {
	qr := "INSERT INTO user(login,username,password) VALUES(?,?,?)"

	_, err := r.db.ExecContext(ctx, qr, user.Login, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	return nil
}

func (r *UserStorage) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	qr := "SELECT id,login,password FROM user WHERE login=$1"

	var user model.User

	row := r.db.QueryRowContext(ctx, qr, login)

	err := row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("couldn't get user: %w", err)
	}

	return user, nil
}

func (r *UserStorage) GetUserByToken(ctx context.Context, token string) (model.User, error) {
	qr := "SELECT * FROM user WHERE token=$1"

	var user model.User

	sqlRow := r.db.QueryRowContext(ctx, qr, token)
	err := sqlRow.Scan(&user.Id, &user.Password, &user.Login, &user.Username, &user.Token, &user.TokenDuration)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserStorage) SaveToken(ctx context.Context, login string, token string, duration time.Time) error {
	qr := "UPDATE user SET token=$1,tokenDuration=$2 WHERE login=$3"

	_, err := r.db.ExecContext(ctx, qr, token, duration, login)
	if err != nil {
		return fmt.Errorf("couldn't update token: %w", err)
	}

	return nil
}

func (r *UserStorage) DeleteToken(ctx context.Context, token string) error {
	qr := "UPDATE user SET token=NULL, tokenDuration=NULL WHERE token=$1"

	_, err := r.db.ExecContext(ctx, qr, token)
	if err != nil {
		return fmt.Errorf("couldn't delete token: %w", err)
	}

	return nil
}
