package mocks

import (
	"errors"
)

type MockFollowService struct {
	Follow      map[int][]int
	UserService MockUserService
}

func NewMockFollowService() *MockFollowService {
	return &MockFollowService{
		Follow:      make(map[int][]int),
		UserService: *NewMockUserService(),
	}
}

func (followService *MockFollowService) Create(followerId int, followeeId int) error {
	if _, err := followService.UserService.GetById(followerId); err != nil {
		if err.Error() == "record not found" {
			return errors.New("ERROR: insert or update on table \"follows\" violates foreign key constraint \"fk_user\"")
		}
		return err
	}
	if _, err := followService.UserService.GetById(followeeId); err != nil {
		if err.Error() == "record not found" {
			return errors.New("ERROR: insert or update on table \"follows\" violates foreign key constraint \"fk_user\"")
		}
		return err
	}
	if followerId == followeeId {
		return errors.New("ERROR: new row for relation \"follows\" violates check constraint \"different_user_and_follower\"")
	}
	if _, ok := followService.Follow[followerId]; ok {
		return errors.New("ERROR: duplicate key value violates unique constraint \"unique_user_follower_pair\"")
	}
	followService.Follow[followerId] = append(followService.Follow[followerId], followeeId)
	return nil
}
