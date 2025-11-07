package unit

import (
	"testing"

	"service/domain/aggregate"
	"service/domain/valueobject"
)

// TestUserRecommendation_CalculateScore 测试推荐分数计算
func TestUserRecommendation_CalculateScore(t *testing.T) {
	// Arrange
	targetUserID, _ := valueobject.NewUserID(100)
	follower1, _ := valueobject.NewUserID(1)
	follower2, _ := valueobject.NewUserID(2)
	follower3, _ := valueobject.NewUserID(3)

	reason := valueobject.NewFollowedByFollowingReason([]valueobject.UserID{
		follower1, follower2, follower3,
	})

	// Act
	recommendation, err := aggregate.NewUserRecommendation(
		targetUserID,
		reason,
		5, // 5个帖子
	)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 预期分数：3个关注者 × 10 + 5个帖子 × 2 = 40
	expectedScore := 40
	if recommendation.Score() != expectedScore {
		t.Errorf("expected score %d, got %d", expectedScore, recommendation.Score())
	}
}

// TestUserRecommendation_NoReason 测试没有推荐理由的情况
func TestUserRecommendation_NoReason(t *testing.T) {
	// Arrange
	targetUserID, _ := valueobject.NewUserID(100)
	reason := valueobject.NewFollowedByFollowingReason([]valueobject.UserID{})

	// Act
	_, err := aggregate.NewUserRecommendation(targetUserID, reason, 0)

	// Assert
	if err == nil {
		t.Error("expected error for recommendation without reason")
	}
}
