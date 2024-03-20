package models

import "time"

type Comment struct {
	Id        int
	CreatedAt time.Time
	PostId    int
	UserId    int
	Content   string
}
