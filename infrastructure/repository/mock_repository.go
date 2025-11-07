package repository

import (
	"context"
	"time"

	"service/application/service"
	"service/domain/entity"
	"service/domain/repository"
	"service/domain/valueobject"
)

// MockSocialGraphRepository Mock 实现：社交图谱仓储
//
// 用于演示和测试，返回模拟数据。
// 在实际项目中，这里会是真实的数据库实现。
type MockSocialGraphRepository struct{}

func NewMockSocialGraphRepository() repository.SocialGraphRepository {
	return &MockSocialGraphRepository{}
}

func (r *MockSocialGraphRepository) GetFollowings(
	ctx context.Context,
	userID valueobject.UserID,
) ([]valueobject.UserID, error) {
	// 返回模拟数据：用户关注了 user2, user3, user4
	user2, _ := valueobject.NewUserID(2)
	user3, _ := valueobject.NewUserID(3)
	user4, _ := valueobject.NewUserID(4)
	return []valueobject.UserID{user2, user3, user4}, nil
}

func (r *MockSocialGraphRepository) GetRecentFollowings(
	ctx context.Context,
	userID valueobject.UserID,
	days int,
) ([]valueobject.UserID, error) {
	// 返回模拟数据：最近关注了 user5, user6
	user5, _ := valueobject.NewUserID(5)
	user6, _ := valueobject.NewUserID(6)
	return []valueobject.UserID{user5, user6}, nil
}

func (r *MockSocialGraphRepository) IsFollowing(
	ctx context.Context,
	followerID, followingID valueobject.UserID,
) (bool, error) {
	// 返回模拟数据：假设存在关注关系
	return true, nil
}

// MockContentRepository Mock 实现：内容仓储
type MockContentRepository struct{}

func NewMockContentRepository() repository.ContentRepository {
	return &MockContentRepository{}
}

func (r *MockContentRepository) CountRecentPosts(
	ctx context.Context,
	userID valueobject.UserID,
	days int,
) (int, error) {
	// 返回模拟数据：5 篇帖子
	return 5, nil
}

func (r *MockContentRepository) GetRecentPosts(
	ctx context.Context,
	userID valueobject.UserID,
	limit int,
) ([]*entity.Post, error) {
	// 返回模拟数据：3 篇帖子
	postID1, _ := valueobject.NewPostID(101)
	postID2, _ := valueobject.NewPostID(102)
	postID3, _ := valueobject.NewPostID(103)

	now := time.Now()
	posts := []*entity.Post{
		entity.NewPost(postID1, userID, "这是第一篇帖子", now.Add(-1*time.Hour)),
		entity.NewPost(postID2, userID, "这是第二篇帖子", now.Add(-2*time.Hour)),
		entity.NewPost(postID3, userID, "这是第三篇帖子", now.Add(-3*time.Hour)),
	}

	return posts, nil
}

// MockUserRPCClient Mock 实现：用户 RPC 客户端
type MockUserRPCClient struct{}

func NewMockUserRPCClient() service.UserRPCClient {
	return &MockUserRPCClient{}
}

func (c *MockUserRPCClient) GetUserInfo(
	ctx context.Context,
	userID int64,
) (*service.UserInfo, error) {
	// 返回模拟数据
	return &service.UserInfo{
		UserID:   userID,
		Username: "user_" + string(rune(userID)),
		Avatar:   "https://example.com/avatar.jpg",
		Bio:      "这是用户简介",
	}, nil
}

func (c *MockUserRPCClient) GetUserInfoBatch(
	ctx context.Context,
	userIDs []int64,
) ([]*service.UserInfo, error) {
	// 返回模拟数据
	result := make([]*service.UserInfo, 0, len(userIDs))
	for _, userID := range userIDs {
		result = append(result, &service.UserInfo{
			UserID:   userID,
			Username: "user_" + string(rune(userID)),
			Avatar:   "https://example.com/avatar.jpg",
			Bio:      "这是用户简介",
		})
	}
	return result, nil
}
