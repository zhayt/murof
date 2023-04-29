package service

import (
	"github.com/gofrs/uuid"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/mail"
	"time"
)

const salt = "Sfasfasfasfas"

var empty model.User

type AuthService struct {
	repo repository.Authorization
	l    *logger.Logger
}

func NewAuthService(repo repository.Authorization, l *logger.Logger) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user model.User) error {
	if !checkLogin(user.Login) {
		return InvalidDate
	}
	if !checkPassword(user.Password) {
		return InvalidDate
	}
	if !checkUsername(user.Username) {
		return InvalidDate
	}

	user.Password = generatePassword(user.Password)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(login, password string) (model.User, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return empty, InvalidDate
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+salt)); err != nil {
		return empty, InvalidDate
	}

	token, err := uuid.NewV4()
	if err != nil {
		log.Print(err)
		return empty, err
	}

	user.Token = token.String()
	user.TokenDuration = time.Now().Add(12 * time.Hour)
	if err := s.repo.SaveTokens(user.Login, user.TokenDuration, user.Token); err != nil {
		return empty, nil
	}
	return user, nil
}

func generatePassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(bytes)
}

func (s *AuthService) GetUserByToken(token string) (model.User, error) {
	return s.repo.GetUserByToken(token)
}

func (s *AuthService) DeleteToken(token string) error {
	return s.repo.DeleteToken(token)
}

func checkUsername(username string) bool {
	for _, w := range username {
		if w >= 33 && w <= 126 {
			return true
		} else {
			return false
		}
	}
	return true
}

func checkPassword(password string) bool {
	for _, rune := range password {
		if rune > 32 && rune <= 126 && len(password) > 8 {
			return true
		} else {
			return false
		}
	}

	return true
}

func checkLogin(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
