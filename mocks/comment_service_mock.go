package mocks

import (
	"errors"
	"math"
	"strconv"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/utils"
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

func (mockCommentService *MockCommentService) ListByPostId(postId int, pageNum string, pageSize string) ([]models.Comment, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}
	var comments []models.Comment
	for _, commentRecord := range mockCommentService.Comments {
		if commentRecord.PostId == postId {
			comments = append(comments, models.Comment{
				Content: commentRecord.Content,
				PostId:  commentRecord.PostId,
				UserId:  commentRecord.UserId,
			})
		}
	}

	totalCount := len(comments)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))

	offset := (pageNumInt - 1) * pageSizeInt
	if offset < 0 {
		offset = 0
	}

	if offset >= totalCount {
		return []models.Comment{}, utils.PageResponse{}, nil
	}

	startIndex := offset
	endIndex := offset + pageSizeInt
	if endIndex > totalCount {
		endIndex = totalCount
	}

	pagedComment := comments[startIndex:endIndex]

	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return pagedComment, pageResponse, nil
}

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
