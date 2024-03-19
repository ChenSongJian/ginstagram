package mocks

import (
	"errors"

	"github.com/ChenSongJian/ginstagram/models"
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

func (userService *MockUserService) GetById(userId int) (models.User, error) {
	for _, user := range userService.Users {
		if user.Id == userId {
			return user, nil
		}
	}
	return models.User{}, errors.New("record not found")
}
