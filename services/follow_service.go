package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type FollowService interface {
	Create(followerId int, followeeId int) error
}

type DBFollowService struct {
	db *gorm.DB
}

func NewDBFollowService() *DBFollowService {
	return &DBFollowService{db: db.DB}
}

func (followService *DBFollowService) Create(followerId int, followeeId int) error {
	result := followService.db.Create(&models.Follow{
		FollowerId: followerId,
		UserId:     followeeId,
	})
	return result.Error
}
