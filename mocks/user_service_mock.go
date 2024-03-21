package mocks

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/utils"
)

type MockUserService struct {
	Users map[string]models.User
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		Users: make(map[string]models.User),
	}
}

func (userService *MockUserService) Create(user models.User) error {
	if _, ok := userService.Users[user.Email]; ok {
		return errors.New("ERROR: duplicate key value violates unique constraint")
	}
	userService.Users[user.Email] = user
	return nil
}

func (userService *MockUserService) List(pageNum string, pageSize string, keyword string) ([]models.User, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	var filteredUsers []models.User
	for _, user := range userService.Users {
		if keyword == "" || strings.Contains(user.Username, keyword) || strings.Contains(user.Bio, keyword) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	totalRecords := len(filteredUsers)
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSizeInt)))

	offset := (pageNumInt - 1) * pageSizeInt
	if offset < 0 {
		offset = 0
	}

	if offset >= totalRecords {
		return []models.User{}, utils.PageResponse{
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
		}, nil
	}

	startIndex := offset
	endIndex := offset + pageSizeInt
	if endIndex > totalRecords {
		endIndex = totalRecords
	}

	pagedUsers := filteredUsers[startIndex:endIndex]

	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	return pagedUsers, pageResponse, nil
}

func (userService *MockUserService) GetById(userId int) (models.User, error) {
	for _, user := range userService.Users {
		if user.Id == userId {
			return user, nil
		}
	}
	return models.User{}, errors.New("record not found")
}

func (userService *MockUserService) GetByEmail(email string) (models.User, error) {
	if user, ok := userService.Users[email]; ok {
		return user, nil
	}
	return models.User{}, errors.New("record not found")
}

func (userService *MockUserService) UpdateByModel(user models.User) error {
	if _, ok := userService.Users[user.Email]; ok {
		userService.Users[user.Email] = user
		return nil
	}
	return errors.New("record not found")
}

func (userService *MockUserService) DeleteById(userId int) error {
	for _, user := range userService.Users {
		if user.Id == userId {
			delete(userService.Users, user.Email)
			return nil
		}
	}
	return errors.New("record not found")
}
