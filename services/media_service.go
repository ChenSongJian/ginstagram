package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type MediaService interface {
	GetByPostId(postId int) ([]models.Media, error)
	Create(media []models.Media) error
}

type DBMediaService struct {
	db *gorm.DB
}

func NewDBMediaService() *DBMediaService {
	return &DBMediaService{db: db.DB}
}

func (mediaService *DBMediaService) Create(media []models.Media) error {
	return mediaService.db.Create(&media).Error
}

func (mediaService *DBMediaService) GetByPostId(postId int) ([]models.Media, error) {
	var media = make([]models.Media, 0)
	err := mediaService.db.Where("post_id = ?", postId).Find(&media).Error
	return media, err
}
