package integration

import (
	"context"
	"testing"

	"service/application/service"
	"service/domain/valueobject"
)

// TestRecommendationService_GetFollowingBasedRecommendations 集成测试
// 使用 mock 仓储测试完整流程
func TestRecommendationService_GetFollowingBasedRecommendations(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// 创建 mock 仓储
	socialGraphRepo := NewMockSocialGraphRepository()
	contentRepo := NewMockContentRepository()
	userRPCClient := NewMockUserRPCClient()

	// 设置测试数据
	// 用户1关注了用户2和用户3
	// 用户2最近关注了用户100
	// 用户3最近关注了用户101
	setupTestData(socialGraphRepo, contentRepo, userRPCClient)

	// 创建服务
	generator := service.NewRecommendationGenerator(socialGraphRepo, contentRepo)
	svc := service.NewRecommendationService(generator, socialGraphRepo, contentRepo, userRPCClient)

	// Act
	result, err := svc.GetFollowingBasedRecommendations(ctx, 1, 10)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Recommendations) != 2 {
		t.Errorf("expected 2 recommendations, got %d", len(result.Recommendations))
	}

	// 验证推荐内容
	for _, rec := range result.Recommendations {
		if rec.UserID != 100 && rec.UserID != 101 {
			t.Errorf("unexpected user id: %d", rec.UserID)
		}
		if rec.Reason == "" {
			t.Error("reason should not be empty")
		}
	}
}

// Mock 实现示例（简化版）
type MockSocialGraphRepository struct {
	// 存储测试数据
}

func NewMockSocialGraphRepository() *MockSocialGraphRepository {
	return &MockSocialGraphRepository{}
}

func (m *MockSocialGraphRepository) GetFollowings(
	ctx context.Context,
	userID valueobject.UserID,
) ([]valueobject.UserID, error) {
	// 返回测试数据
	return []valueobject.UserID{}, nil
}

// ... 其他 mock 实现
