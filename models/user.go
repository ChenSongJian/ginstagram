package models

import "time"

type User struct {
	Id              int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Username        string
	PasswordHash    string
	Email           string
	Bio             string
	ProfileImageUrl string
}
