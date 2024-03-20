package services

import (
	"fmt"

	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type PostService interface {
	Create(post models.Post) (int, error)
	DeleteById(id int) error
}

type DBPostService struct {
	db *gorm.DB
}

func NewDBPostService() *DBPostService {
	return &DBPostService{db: db.DB}
}

func (postService *DBPostService) Create(post models.Post) (int, error) {
	result := postService.db.Create(&post)
	if result.Error != nil {
		return 0, result.Error
	}
	fmt.Println(post.Id)
	return post.Id, nil
}

func (postService *DBPostService) DeleteById(id int) error {
	result := postService.db.Delete(&models.Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
