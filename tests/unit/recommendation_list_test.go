package unit

import (
	"testing"

	"service/domain/aggregate"
	"service/domain/valueobject"
)

// TestRecommendationList_CannotRecommendSelf 测试不能推荐自己
func TestRecommendationList_CannotRecommendSelf(t *testing.T) {
	// Arrange
	userID, _ := valueobject.NewUserID(1)
	list := aggregate.NewRecommendationList(userID)

	follower, _ := valueobject.NewUserID(2)
	reason := valueobject.NewFollowedByFollowingReason([]valueobject.UserID{follower})

	// 创建推荐自己的推荐
	recommendation, _ := aggregate.NewUserRecommendation(userID, reason, 0)

	// Act
	err := list.AddRecommendation(recommendation)

	// Assert
	if err == nil {
		t.Error("expected error when recommending self")
	}
}

// TestRecommendationList_NoDuplicates 测试不能重复推荐
func TestRecommendationList_NoDuplicates(t *testing.T) {
	// Arrange
	userID, _ := valueobject.NewUserID(1)
	targetUserID, _ := valueobject.NewUserID(100)
	list := aggregate.NewRecommendationList(userID)

	follower, _ := valueobject.NewUserID(2)
	reason := valueobject.NewFollowedByFollowingReason([]valueobject.UserID{follower})

	rec1, _ := aggregate.NewUserRecommendation(targetUserID, reason, 0)
	rec2, _ := aggregate.NewUserRecommendation(targetUserID, reason, 0)

	// Act
	err1 := list.AddRecommendation(rec1)
	err2 := list.AddRecommendation(rec2)

	// Assert
	if err1 != nil {
		t.Errorf("first add should succeed: %v", err1)
	}
	if err2 == nil {
		t.Error("second add should fail (duplicate)")
	}
}
