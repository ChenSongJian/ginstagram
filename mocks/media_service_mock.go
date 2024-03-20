package mocks

import "github.com/ChenSongJian/ginstagram/models"

type MockMediaService struct {
	Media map[int]MediaRecord
}

func NewMockMediaService() *MockMediaService {
	return &MockMediaService{
		Media: map[int]MediaRecord{},
	}
}

type MediaRecord struct {
	Url    string
	PostId int
}

var MediaRecordId = 0

func (mediaService *MockMediaService) Create(media []models.Media) error {
	for _, m := range media {
		MediaRecordId++
		mediaService.Media[MediaRecordId] = MediaRecord{
			Url:    m.Url,
			PostId: m.PostId,
		}
	}
	return nil
}

func (mediaService *MockMediaService) DeleteByPostId(postId int) error {
	var mediaId int
	for k, v := range mediaService.Media {
		if v.PostId == postId {
			mediaId = k
			break
		}
	}
	delete(mediaService.Media, mediaId)
	return nil
}
