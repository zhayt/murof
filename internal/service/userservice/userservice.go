package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userManipulation interface {
	CreateUser(ctx context.Context, user *model.User) error
	User(ctx context.Context, email string) (*model.User, error)
	CrateSession(ctx context.Context, user *model.User) error
	Session(ctx context.Context, userId int) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
}

type UserService struct {
	userManipulation userManipulation
}

func NewUserService(manipulation userManipulation) *UserService {
	return &UserService{userManipulation: manipulation}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	hash, err := generatePasswordHash(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	user.PasswordHash = hash

	if err := s.userManipulation.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	return nil
}

func (s *UserService) User(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.userManipulation.User(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can't get user: %w", err)
	}

	if err = compareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("passwordhash and pasword not equal: %w", err)
	}

	session, err := s.userManipulation.Session(ctx, user.ID)
	if err == sql.ErrNoRows {
		token, err := generateToken()
		if err != nil {
			return nil, fmt.Errorf("can't generate token: %w)", err)
		}

		user.Token = token
		user.ExpirationTime = time.Now().Add(12 * time.Hour)

		if err = s.userManipulation.CrateSession(ctx, user); err != nil {
			return nil, fmt.Errorf("can't create session: %w", err)
		}

		return user, nil
	}

	if err != nil {
		return nil, fmt.Errorf("can't get session: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("can't generate token: %w)", err)
	}

	session.Token = token
	session.ExpirationTime = time.Now().Add(12 * time.Hour)
	if err = s.userManipulation.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("can't update session: %w", err)
	}

	user.Token = session.Token
	user.ExpirationTime = session.ExpirationTime

	return user, nil
}

func generatePasswordHash(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate password hash: %w", err)
	}

	return string(hash), nil
}

func compareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func generateToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
