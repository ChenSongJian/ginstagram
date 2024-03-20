package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type CommentService interface {
	Create(postId int, userId int, content string) error
}

type DBCommentService struct {
	db *gorm.DB
}

func NewDBCommentService() *DBCommentService {
	return &DBCommentService{db: db.DB}
}

func (commentService *DBCommentService) Create(postId int, userId int, content string) error {
	return commentService.db.Create(&models.Comment{
		PostId:  postId,
		UserId:  userId,
		Content: content,
	}).Error
}
