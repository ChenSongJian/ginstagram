package mocks

import (
	"errors"
	"math"
	"strconv"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/utils"
)

type MockPostService struct {
	Posts         map[int]PostRecord
	UserService   *MockUserService
	FollowService *MockFollowService
	MediaService  *MockMediaService
}

func NewMockPostService() *MockPostService {
	return &MockPostService{
		Posts:         map[int]PostRecord{},
		UserService:   NewMockUserService(),
		FollowService: NewMockFollowService(),
		MediaService:  NewMockMediaService(),
	}
}

type PostRecord struct {
	Title   string
	Content string
	UserId  int
}

var PostRecordId = 0

func (postService *MockPostService) List(pageNum string, pageSize string, keyword string) ([]models.Post, map[int][]models.Media, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	filterUserIds := []int{}
	for _, user := range postService.UserService.Users {
		if !user.IsPrivate {
			filterUserIds = append(filterUserIds, user.Id)
		}
	}

	var posts []models.Post
	var postIds []int
	if keyword != "" {
		for id, post := range postService.Posts {
			if (post.Title == keyword || post.Content == keyword) && utils.IsInIntSlice(post.UserId, filterUserIds) {
				posts = append(posts, models.Post{
					Id:      id,
					Title:   post.Title,
					Content: post.Content,
					UserId:  post.UserId,
				})
				postIds = append(postIds, id)
			}
		}
	} else {
		for id, post := range postService.Posts {
			if utils.IsInIntSlice(post.UserId, filterUserIds) {
				posts = append(posts, models.Post{
					Id:      id,
					Title:   post.Title,
					Content: post.Content,
					UserId:  post.UserId,
				})
				postIds = append(postIds, id)
			}
		}
	}
	mediaMap := make(map[int][]models.Media)
	if len(postIds) > 0 {
		for _, m := range postService.MediaService.Media {
			if utils.IsInIntSlice(m.PostId, postIds) {
				mediaMap[m.PostId] = append(mediaMap[m.PostId], models.Media{
					Url: m.Url,
				})
			}
		}
	}
	totalCount := len(posts)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))

	offset := (pageNumInt - 1) * pageSizeInt
	if offset < 0 {
		offset = 0
	}

	if offset >= totalCount {
		return []models.Post{}, map[int][]models.Media{}, utils.PageResponse{
			TotalPages:   totalPages,
			TotalRecords: totalCount,
		}, nil
	}

	startIndex := offset
	endIndex := offset + pageSizeInt
	if endIndex > totalCount {
		endIndex = totalCount
	}

	pagedPost := posts[startIndex:endIndex]

	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return pagedPost, mediaMap, pageResponse, nil
}

func (postService *MockPostService) ListByUserId(userId int, pageNum string, pageSize string, keyword string) ([]models.Post, map[int][]models.Media, utils.PageResponse, error) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	filterUserIds := []int{userId}
	for _, user := range postService.UserService.Users {
		if user.Id != userId && !user.IsPrivate {
			filterUserIds = append(filterUserIds, user.Id)
		}
	}
	for _, follow := range postService.FollowService.Follows {
		if follow.FollowerId == userId {
			filterUserIds = append(filterUserIds, follow.FolloweeId)
		}
	}

	var posts []models.Post
	var postIds []int
	if keyword != "" {
		for id, post := range postService.Posts {
			if (post.Title == keyword || post.Content == keyword) && utils.IsInIntSlice(post.UserId, filterUserIds) {
				posts = append(posts, models.Post{
					Id:      id,
					Title:   post.Title,
					Content: post.Content,
					UserId:  post.UserId,
				})
				postIds = append(postIds, id)
			}
		}
	} else {
		for id, post := range postService.Posts {
			if utils.IsInIntSlice(post.UserId, filterUserIds) {
				posts = append(posts, models.Post{
					Id:      id,
					Title:   post.Title,
					Content: post.Content,
					UserId:  post.UserId,
				})
				postIds = append(postIds, id)
			}
		}
	}
	mediaMap := make(map[int][]models.Media)
	if len(postIds) > 0 {
		for _, m := range postService.MediaService.Media {
			if utils.IsInIntSlice(m.PostId, postIds) {
				mediaMap[m.PostId] = append(mediaMap[m.PostId], models.Media{
					Url: m.Url,
				})
			}
		}
	}
	totalCount := len(posts)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))

	offset := (pageNumInt - 1) * pageSizeInt
	if offset < 0 {
		offset = 0
	}

	if offset >= totalCount {
		return []models.Post{}, map[int][]models.Media{}, utils.PageResponse{
			TotalPages:   totalPages,
			TotalRecords: totalCount,
		}, nil
	}

	startIndex := offset
	endIndex := offset + pageSizeInt
	if endIndex > totalCount {
		endIndex = totalCount
	}

	pagedPost := posts[startIndex:endIndex]

	pageResponse := utils.PageResponse{
		PageNum:      pageNumInt,
		PageSize:     pageSizeInt,
		TotalPages:   totalPages,
		TotalRecords: int(totalCount),
	}
	return pagedPost, mediaMap, pageResponse, nil
}

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
