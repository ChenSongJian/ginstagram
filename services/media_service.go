package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type MediaService interface {
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
