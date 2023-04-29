package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"log"
	"time"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) CreateUser(user model.User) error {
	query := "INSERT INTO user(login,username,password) VALUES(?,?,?)"
	_, err := a.db.Exec(query, user.Login, user.Username, user.Password)
	if err != nil {
		log.Fatal("create user repo", err)
	}
	return nil
}

func (a *AuthRepo) GetUser(login string) (model.User, error) {
	var fullUser model.User
	query := "SELECT id,login,password FROM user WHERE login=$1"
	sqlRow := a.db.QueryRow(query, login)
	err := sqlRow.Scan(&fullUser.Id, &fullUser.Login, &fullUser.Password)
	if err != nil {
		return fullUser, err
	}
	return fullUser, nil
}

func (a *AuthRepo) GetUserByToken(token string) (model.User, error) {
	var fullUser model.User
	query := "SELECT * FROM user WHERE token=$1"
	sqlRow := a.db.QueryRow(query, token)
	err := sqlRow.Scan(&fullUser.Id, &fullUser.Password, &fullUser.Login, &fullUser.Username, &fullUser.Token, &fullUser.TokenDuration)
	if err != nil {
		return fullUser, err
	}
	return fullUser, nil
}

func (a *AuthRepo) SaveTokens(login string, tokenDuration time.Time, token string) error {
	query := "UPDATE user SET token=$1,tokenDuration=$2 WHERE login=$3"
	_, err := a.db.Exec(query, token, tokenDuration, login)
	if err != nil {
		return fmt.Errorf("ERROR:don't save user's token: %w", err)
	}
	return nil
}

func (a *AuthRepo) DeleteToken(token string) error {
	query := "UPDATE user SET token=NULL, tokenDuration=NULL WHERE token=$1"
	_, err := a.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("ERROR: don't delete user's token: %w", err)
	}
	return nil
}
