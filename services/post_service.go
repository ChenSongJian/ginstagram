package services

import (
	"math"
	"strconv"

	"github.com/ChenSongJian/ginstagram/db"
	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/utils"
	"gorm.io/gorm"
)

type PostService interface {
	List(pageNum string, pageSize string, keyword string) ([]models.Post, map[int][]models.Media, utils.PageResponse, error)
	ListByUserId(userId int, pageNum string, pageSize string, keyword string) ([]models.Post, map[int][]models.Media, utils.PageResponse, error)
	GetById(postId int) (models.Post, error)
	Create(post models.Post) (int, error)
	DeleteById(id int) error
}

type DBPostService struct {
	db *gorm.DB
}

func NewDBPostService() *DBPostService {
	return &DBPostService{db: db.DB}
}

func (postService *DBPostService) List(pageNum string, pageSize string, keyword string) ([]models.Post, map[int][]models.Media, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}
	offset := (pageNumInt - 1) * pageSizeInt

	filterUserIds := []int{}
	var publicUsers []models.User
	postService.db.Where("is_private=false").Find(&publicUsers)
	for _, publicUser := range publicUsers {
		filterUserIds = append(filterUserIds, publicUser.Id)
	}

	var posts []models.Post
	query := postService.db.Model(&models.Post{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query = query.Where("user_id IN ?", filterUserIds).Find(&posts)
	var totalCount int64
	query.Count(&totalCount)
	query.Order("created_at desc").Offset(offset).Limit(pageSizeInt).Find(&posts)
	var postIds []int
	for _, post := range posts {
		postIds = append(postIds, post.Id)
	}
	mediaMap := make(map[int][]models.Media)
	if len(postIds) > 0 {
		var media []models.Media
		postService.db.Where("post_id IN ?", postIds).Find(&media)
		for _, m := range media {
			mediaMap[m.PostId] = append(mediaMap[m.PostId], m)
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))
	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return posts, mediaMap, pageResponse, nil
}

func (postService *DBPostService) ListByUserId(userId int, pageNum string, pageSize string, keyword string) ([]models.Post, map[int][]models.Media, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}
	offset := (pageNumInt - 1) * pageSizeInt

	filterUserIds := []int{userId}
	var followingUsers []models.Follow
	postService.db.Where("follower_id = ?", userId).Find(&followingUsers)
	for _, followingUser := range followingUsers {
		filterUserIds = append(filterUserIds, followingUser.UserId)
	}
	var publicUsers []models.User
	postService.db.Where("is_private=false").Find(&publicUsers)
	for _, publicUser := range publicUsers {
		filterUserIds = append(filterUserIds, publicUser.Id)
	}

	var posts []models.Post
	query := postService.db.Model(&models.Post{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query = query.Where("user_id IN ?", filterUserIds).Find(&posts)
	var totalCount int64
	query.Count(&totalCount)
	query.Order("created_at desc").Offset(offset).Limit(pageSizeInt).Find(&posts)
	var postIds []int
	for _, post := range posts {
		postIds = append(postIds, post.Id)
	}
	mediaMap := make(map[int][]models.Media)
	if len(postIds) > 0 {
		var media []models.Media
		postService.db.Where("post_id IN ?", postIds).Find(&media)
		for _, m := range media {
			mediaMap[m.PostId] = append(mediaMap[m.PostId], m)
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))
	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return posts, mediaMap, pageResponse, nil
}

func (postService *DBPostService) GetById(postId int) (models.Post, error) {
	var post models.Post
	result := postService.db.First(&post, postId)
	if result.Error != nil {
		return post, result.Error
	}
	return post, nil
}

func (postService *DBPostService) Create(post models.Post) (int, error) {
	result := postService.db.Create(&post)
	if result.Error != nil {
		return 0, result.Error
	}
	return post.Id, nil
}

func (postService *DBPostService) DeleteById(id int) error {
	result := postService.db.Delete(&models.Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
