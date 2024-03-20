package services

import (
	"fmt"

	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"gorm.io/gorm"
)

type PostService interface {
	GetById(postId int) (models.Post, error)
	Create(post models.Post) (int, error)
	DeleteById(id int) error
}

type DBPostService struct {
	db *gorm.DB
}

func NewDBPostService() *DBPostService {
	return &DBPostService{db: db.DB}
}

func (postService *DBPostService) GetById(postId int) (models.Post, error) {
	var post models.Post
	result := postService.db.First(&post, postId)
	if result.Error != nil {
		return post, result.Error
	}
	return post, nil
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
