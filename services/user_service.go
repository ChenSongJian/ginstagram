package services

import (
	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user models.User) error
	GetById(userId int) (models.User, error)
}

type DBUserService struct {
	db *gorm.DB
}

func NewDBUserService() *DBUserService {
	return &DBUserService{db: db.DB}
}

func (userService DBUserService) Create(user models.User) error {
	result := userService.db.Create(&user)
	return result.Error
}

func (userService DBUserService) GetById(userId int) (models.User, error) {
	var user models.User
	result := userService.db.First(&user, userId)
	return user, result.Error
}
