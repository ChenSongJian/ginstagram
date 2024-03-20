package services

import (
	"math"
	"strconv"

	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/utils"
	"gorm.io/gorm"
)

type CommentService interface {
	ListByPostId(postId int, pageNum string, pageSize string) ([]models.Comment, utils.PageResponse, error)
	GetById(commentId int) (models.Comment, error)
	Create(postId int, userId int, content string) error
	DeleteById(commentId int) error
}

type DBCommentService struct {
	db *gorm.DB
}

func NewDBCommentService() *DBCommentService {
	return &DBCommentService{db: db.DB}
}

func (commentService *DBCommentService) ListByPostId(postId int, pageNum string, pageSize string) ([]models.Comment, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}
	offset := (pageNumInt - 1) * pageSizeInt

	var comments []models.Comment
	query := commentService.db.Where("post_id = ?", postId).Order("created_at desc").Offset(offset).Limit(pageSizeInt)
	query.Find(&comments)

	var totalCount int64
	query.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))
	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return comments, pageResponse, nil
}

func (commentService *DBCommentService) Create(postId int, userId int, content string) error {
	return commentService.db.Create(&models.Comment{
		PostId:  postId,
		UserId:  userId,
		Content: content,
	}).Error
}

func (commentService *DBCommentService) GetById(commentId int) (models.Comment, error) {
	var comment models.Comment
	err := commentService.db.First(&comment, commentId).Error
	return comment, err
}

func (commentService *DBCommentService) DeleteById(commentId int) error {
	return commentService.db.Delete(&models.Comment{}, commentId).Error
}
