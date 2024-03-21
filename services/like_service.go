package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type LikeService interface {
	GetByPostLikeId(postLikeId int) (models.PostLike, error)
	CreatePostLike(postId int, userId int) error
	DeletePostLikeById(postLikeId int) error
}

type DBLikeService struct {
	db *gorm.DB
}

func NewDBLikeService() *DBLikeService {
	return &DBLikeService{db: db.DB}
}

func (likeService *DBLikeService) GetByPostLikeId(postLikeId int) (models.PostLike, error) {
	var postLike models.PostLike
	err := likeService.db.First(postLike, postLikeId).Error
	return postLike, err
}

func (likeService *DBLikeService) CreatePostLike(postId int, userId int) error {
	return likeService.db.Create(&models.PostLike{PostId: postId, UserId: userId}).Error
}

func (likeService *DBLikeService) DeletePostLikeById(postLikeId int) error {
	return likeService.db.Delete(&models.PostLike{}, postLikeId).Error
}
