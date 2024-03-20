package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type CommentService interface {
	GetById(commentId int) (models.Comment, error)
	Create(postId int, userId int, content string) error
	DeleteById(commentId int) error
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

func (commentService *DBCommentService) GetById(commentId int) (models.Comment, error) {
	var comment models.Comment
	err := commentService.db.First(&comment, commentId).Error
	return comment, err
}

func (commentService *DBCommentService) DeleteById(commentId int) error {
	return commentService.db.Delete(&models.Comment{}, commentId).Error
}
