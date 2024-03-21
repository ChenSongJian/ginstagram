package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type LikeService interface {
	ListPostLikesByPostId(postId int) ([]models.PostLike, error)
	GetByPostLikeId(postLikeId int) (models.PostLike, error)
	CreatePostLike(postId int, userId int) error
	DeletePostLikeById(postLikeId int) error

	GetByCommentLikeId(commentId int) (models.CommentLike, error)
	CreateCommentLike(commentId int, userId int) error
	DeleteCommentLikeById(commentLikeId int) error
}

type DBLikeService struct {
	db *gorm.DB
}

func NewDBLikeService() *DBLikeService {
	return &DBLikeService{db: db.DB}
}

func (likeService *DBLikeService) ListPostLikesByPostId(postId int) ([]models.PostLike, error) {
	var postLikes []models.PostLike
	err := likeService.db.Where("post_id = ?", postId).Find(&postLikes).Error
	return postLikes, err
}

func (likeService *DBLikeService) GetByPostLikeId(postLikeId int) (models.PostLike, error) {
	var postLike models.PostLike
	err := likeService.db.First(&postLike, postLikeId).Error
	return postLike, err
}

func (likeService *DBLikeService) CreatePostLike(postId int, userId int) error {
	return likeService.db.Create(&models.PostLike{PostId: postId, UserId: userId}).Error
}

func (likeService *DBLikeService) DeletePostLikeById(postLikeId int) error {
	return likeService.db.Delete(&models.PostLike{}, postLikeId).Error
}

func (likeService *DBLikeService) GetByCommentLikeId(commentId int) (models.CommentLike, error) {
	var commentLike models.CommentLike
	err := likeService.db.First(&commentLike, commentId).Error
	return commentLike, err
}

func (likeService *DBLikeService) CreateCommentLike(commentId int, userId int) error {
	return likeService.db.Create(&models.CommentLike{CommentId: commentId, UserId: userId}).Error
}

func (likeService *DBLikeService) DeleteCommentLikeById(commentLikeId int) error {
	return likeService.db.Delete(&models.CommentLike{}, commentLikeId).Error
}
