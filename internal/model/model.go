package model

import "time"

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string

	Token          string
	ExpirationTime time.Time
}

type Session struct {
	ID             int
	UserId         int
	Token          string
	ExpirationTime time.Time
}
