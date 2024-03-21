package models

type PostLike struct {
	Id     int
	UserId int
	PostId int
}

type CommentLike struct {
	Id        int
	UserId    int
	CommentId int
}
