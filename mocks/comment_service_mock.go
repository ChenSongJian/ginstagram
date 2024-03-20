package mocks

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
