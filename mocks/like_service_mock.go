package mocks

import (
	"errors"

	"github.com/ChenSongJian/ginstagram/models"
)

type MockLikeService struct {
	PostLikes     map[int]PostLikeRecord
	CommentLikes  map[int]CommentLikeRecord
	UserService   *MockUserService
	FollowService *MockFollowService
	PostService   *MockPostService
}

func NewMockLikeService() *MockLikeService {
	return &MockLikeService{
		PostLikes:     make(map[int]PostLikeRecord),
		CommentLikes:  make(map[int]CommentLikeRecord),
		UserService:   NewMockUserService(),
		FollowService: NewMockFollowService(),
		PostService:   NewMockPostService(),
	}
}

type PostLikeRecord struct {
	UserId int
	PostId int
}
type CommentLikeRecord struct {
	UserId    int
	CommentId int
}

var PostLikeRecordId = 0
var CommentLikeRecordId = 0

func (likeService *MockLikeService) GetByPostLikeId(postLikeId int) (models.PostLike, error) {
	if like, ok := likeService.PostLikes[postLikeId]; ok {
		return models.PostLike{
			UserId: like.UserId,
			PostId: like.PostId,
		}, nil
	}
	return models.PostLike{}, errors.New("record not found")
}

func (likeService *MockLikeService) CreatePostLike(userId int, postId int) error {
	userFound := false
	for _, user := range likeService.UserService.Users {
		if user.Id == userId {
			userFound = true
		}
	}
	if !userFound {
		return errors.New("violates foreign key constraint \"post_likes_user_id_fkey\"")
	}

	if _, ok := likeService.PostService.Posts[postId]; !ok {
		return errors.New("violates foreign key constraint \"post_likes_post_id_fkey\"")
	}

	for _, like := range likeService.PostLikes {
		if like.UserId == userId && like.PostId == postId {
			return errors.New("violates unique constraint \"unique_post_user_pair\"")
		}
	}

	PostLikeRecordId++
	likeService.PostLikes[PostLikeRecordId] = PostLikeRecord{
		UserId: userId,
		PostId: postId,
	}
	return nil
}

func (likeService *MockLikeService) DeletePostLikeById(postLikeId int) error {
	delete(likeService.PostLikes, postLikeId)
	return nil
}
