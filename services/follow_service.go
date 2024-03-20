package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type FollowService interface {
	GetById(followId int) (models.Follow, error)
	GetByFollowerId(followId int) ([]models.Follow, error)
	GetByFolloweeId(followId int) ([]models.Follow, error)
	Create(followerId int, followeeId int) error
	Delete(followId int) error
}

type DBFollowService struct {
	db *gorm.DB
}

func NewDBFollowService() *DBFollowService {
	return &DBFollowService{db: db.DB}
}

func (followService *DBFollowService) GetById(followId int) (models.Follow, error) {
	var follow models.Follow
	result := followService.db.First(&follow, followId)
	return follow, result.Error
}

func (followService *DBFollowService) GetByFollowerId(followerId int) ([]models.Follow, error) {
	var follows []models.Follow
	result := followService.db.Where("follower_id = ?", followerId).Find(&follows)
	return follows, result.Error
}

func (followService *DBFollowService) GetByFolloweeId(followeeId int) ([]models.Follow, error) {
	var follows []models.Follow
	result := followService.db.Where("user_id = ?", followeeId).Find(&follows)
	return follows, result.Error
}

func (followService *DBFollowService) Create(followerId int, followeeId int) error {
	result := followService.db.Create(&models.Follow{
		FollowerId: followerId,
		UserId:     followeeId,
	})
	return result.Error
}

func (followService *DBFollowService) Delete(followId int) error {
	result := followService.db.Delete(&models.Follow{}, followId)
	return result.Error
}
