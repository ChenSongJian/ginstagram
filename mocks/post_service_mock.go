package mocks

import (
	"errors"

	"github.com/ChenSongJian/ginstagram/models"
)

type MockPostService struct {
	Posts        map[int]PostRecord
	UserService  *MockUserService
	MediaService *MockMediaService
}

func NewMockPostService() *MockPostService {
	return &MockPostService{
		Posts:        map[int]PostRecord{},
		UserService:  NewMockUserService(),
		MediaService: NewMockMediaService(),
	}
}

type PostRecord struct {
	Title   string
	Content string
	UserId  int
}

var PostRecordId = 0

func (postService *MockPostService) GetById(postId int) (models.Post, error) {
	postRecord, ok := postService.Posts[postId]
	if !ok {
		return models.Post{}, errors.New("record not found")
	}
	post := models.Post{
		Id:      postId,
		Title:   postRecord.Title,
		Content: postRecord.Content,
		UserId:  postRecord.UserId,
	}
	return post, nil
}

func (postService *MockPostService) Create(post models.Post) (int, error) {
	PostRecordId++
	postRecord := PostRecord{
		Title:   post.Title,
		Content: post.Content,
		UserId:  post.UserId,
	}
	postService.Posts[PostRecordId] = postRecord
	return PostRecordId, nil
}

func (postService *MockPostService) DeleteById(postId int) error {
	delete(postService.Posts, postId)
	postService.MediaService.DeleteByPostId(postId)
	return nil
}
