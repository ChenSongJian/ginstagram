package mocks

import (
	"errors"

	"github.com/ChenSongJian/ginstagram/models"
)

type MockCommentService struct {
	Comments      map[int]CommentRecord
	UserService   *MockUserService
	FollowService *MockFollowService
	PostService   *MockPostService
}

func NewMockCommentService() *MockCommentService {
	return &MockCommentService{
		Comments:      make(map[int]CommentRecord),
		UserService:   NewMockUserService(),
		FollowService: NewMockFollowService(),
		PostService:   NewMockPostService(),
	}
}

type CommentRecord struct {
	Content string
	PostId  int
	UserId  int
}

var commentRecordId = 0

func (mockCommentService *MockCommentService) Create(postId int, userId int, content string) error {
	commentRecordId++
	commentRecord := CommentRecord{
		Content: content,
		PostId:  postId,
		UserId:  userId,
	}
	mockCommentService.Comments[commentRecordId] = commentRecord
	return nil
}

func (mockCommentService *MockCommentService) GetById(commentId int) (models.Comment, error) {
	commentRecord, ok := mockCommentService.Comments[commentId]
	if !ok {
		return models.Comment{}, errors.New("record not found")
	}
	return models.Comment{
		Content: commentRecord.Content,
		PostId:  commentRecord.PostId,
		UserId:  commentRecord.UserId,
	}, nil
}

func (mockCommentService *MockCommentService) DeleteById(commentId int) error {
	delete(mockCommentService.Comments, commentId)
	return nil
}
