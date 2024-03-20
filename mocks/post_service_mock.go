package mocks

import "github.com/ChenSongJian/ginstagram/models"

type MockPostService struct {
	Post        map[int]PostRecord
	userService *MockUserService
}

func NewMockPostService() *MockPostService {
	return &MockPostService{
		Post:        map[int]PostRecord{},
		userService: NewMockUserService(),
	}
}

type PostRecord struct {
	title   string
	content string
	userId  int
}

var PostRecordId = 0

func (postService *MockPostService) Create(post models.Post) (int, error) {
	PostRecordId++
	postRecord := PostRecord{
		title:   post.Title,
		content: post.Title,
		userId:  post.UserId,
	}
	postService.Post[PostRecordId] = postRecord
	return PostRecordId, nil
}

func (postService *MockPostService) DeleteById(id int) error {
	delete(postService.Post, id)
	return nil
}
