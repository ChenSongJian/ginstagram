package services

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/utils"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user models.User) error
	List(pageNum string, pageSize string, keyword string) ([]models.User, utils.PageResponse, error)
	GetById(userId int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	UpdateByModel(user models.User) error
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

func (userService DBUserService) List(pageNum string, pageSize string, keyword string) ([]models.User, utils.PageResponse, error) {
	var users []models.User
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}
	offset := (pageNumInt - 1) * pageSizeInt

	var totalCount int64
	query := userService.db.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR bio LIKE ?", fmt.Sprintf("%%%s%%", keyword), fmt.Sprintf("%%%s%%", keyword))
	}
	query.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))
	if err := query.Limit(pageSizeInt).Offset(offset).Find(&users).Error; err != nil {
		return users, utils.PageResponse{}, err
	}

	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return users, pageResponse, nil
}

func (userService DBUserService) GetById(userId int) (models.User, error) {
	var user models.User
	result := userService.db.First(&user, userId)
	return user, result.Error
}

func (userService DBUserService) GetByEmail(email string) (models.User, error) {
	var user models.User
	result := userService.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (userService DBUserService) UpdateByModel(modelUser models.User) error {
	var user models.User
	result := userService.db.First(&user, modelUser.Id)
	if result.Error != nil {
		return result.Error
	}
	result = userService.db.Save(&modelUser)
	return result.Error
}
