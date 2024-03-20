package mocks

import (
	"errors"

	"github.com/ChenSongJian/ginstagram/models"
)

type MockFollowService struct {
	Follows     map[int]FollowRecord
	UserService MockUserService
}

func NewMockFollowService() *MockFollowService {
	return &MockFollowService{
		Follows:     make(map[int]FollowRecord),
		UserService: *NewMockUserService(),
	}
}

type FollowRecord struct {
	FollowerId int
	FolloweeId int
}

var followRecordId = 0

func (followService *MockFollowService) GetById(followId int) (models.Follow, error) {
	if follow, ok := followService.Follows[followId]; ok {
		return models.Follow{
			Id:         followId,
			FollowerId: follow.FollowerId,
			UserId:     follow.FolloweeId,
		}, nil
	}
	return models.Follow{}, errors.New("record not found")
}

func (followService *MockFollowService) GetByFollowerId(followerId int) ([]models.Follow, error) {
	var result = make([]models.Follow, 0)
	for k, v := range followService.Follows {
		if v.FollowerId == followerId {
			result = append(result, models.Follow{
				Id:         k,
				FollowerId: v.FollowerId,
				UserId:     v.FolloweeId,
			})
		}
	}
	return result, nil
}

func (followService *MockFollowService) GetByFolloweeId(followeeId int) ([]models.Follow, error) {
	var result = make([]models.Follow, 0)
	for k, v := range followService.Follows {
		if v.FolloweeId == followeeId {
			result = append(result, models.Follow{
				Id:         k,
				FollowerId: v.FollowerId,
				UserId:     v.FolloweeId,
			})
		}
	}
	return result, nil
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

	record := FollowRecord{
		FollowerId: followerId,
		FolloweeId: followeeId,
	}

	for _, v := range followService.Follows {
		if v == record {
			return errors.New("ERROR: duplicate key value violates unique constraint \"follows_pkey\"")
		}
	}
	followRecordId++
	followService.Follows[followRecordId] = record
	return nil
}

func (followService *MockFollowService) Delete(followId int) error {
	delete(followService.Follows, followId)
	return nil
}

func (followService *MockFollowService) IsFollowing(followerId int, followeeId int) bool {
	for _, v := range followService.Follows {
		if v.FollowerId == followerId && v.FolloweeId == followeeId {
			return true
		}
	}
	return false
}
