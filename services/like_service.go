package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type LikeService interface {
	CreatePostLike(postId int, userId int) error
}

type DBLikeService struct {
	db *gorm.DB
}

func NewDBLikeService() *DBLikeService {
	return &DBLikeService{db: db.DB}
}

func (likeService *DBLikeService) CreatePostLike(postId int, userId int) error {
	return likeService.db.Create(&models.PostLike{PostId: postId, UserId: userId}).Error
}
