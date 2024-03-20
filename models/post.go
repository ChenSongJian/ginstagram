package models

import "time"

type Post struct {
	Id        int
	CreatedAt time.Time
	Title     string
	Content   string
	UserId    int
}
