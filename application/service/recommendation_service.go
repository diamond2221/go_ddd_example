package service

import (
	"context"

	"service/application/dto"
	"service/domain/repository"
	"service/domain/service"

	"service/domain/entity"
	"service/domain/valueobject"
)

// RecommendationService 应用服务：推荐用例编排
//
// 什么是应用服务？
// 应用服务是应用层的核心，负责编排用例（Use Case）。
// 它协调领域对象、基础设施服务来完成一个完整的业务流程。
//
// 应用服务的职责（记住：编排，不是实现）：
// 1. 用例编排：协调领域对象完成业务用例
// 2. 跨服务调用：调用其他微服务（如 user 服务）
// 3. DTO 转换：领域对象 ↔ 数据传输对象
// 4. 事务管理：控制事务边界
// 5. 权限检查：验证用户是否有权限执行操作
//
// 应用服务不应该包含：
// - 业务规则：应该在领域层
// - 数据库访问：应该通过仓储
// - 协议细节：应该在接口层
//
// 应用服务 vs 领域服务：
// ┌──────────────────────────────────────────────────────┐
// │ 场景：生成推荐并返回给用户                            │
// ├──────────────────────────────────────────────────────┤
// │ 领域服务（RecommendationGenerator）：                 │
// │   - 实现推荐算法                                      │
// │   - 创建推荐聚合                                      │
// │   - 纯业务逻辑                                        │
// ├──────────────────────────────────────────────────────┤
// │ 应用服务（RecommendationService）：                   │
// │   1. 调用领域服务生成推荐                             │
// │   2. 调用 user 服务获取用户信息（跨服务）             │
// │   3. 调用 content 服务获取帖子（跨服务）              │
// │   4. 组装数据并转换为 DTO                             │
// │   5. 返回给接口层                                     │
// └──────────────────────────────────────────────────────┘
//
// 实际业务场景：
// 用户点击"推荐"按钮 →
//
//	接口层接收请求 →
//	  应用服务编排用例 →
//	    领域服务生成推荐 →
//	    RPC 获取用户信息 →
//	    组装响应数据 →
//	  返回给接口层 →
//	返回给用户
//
// 对比传统方式：
// 传统方式：所有逻辑都在 Service 层，业务规则和技术细节混在一起
// DDD 方式：业务规则在领域层，应用服务只负责编排
type RecommendationService struct {
	generator          *service.RecommendationGenerator
	socialGraphRepo    repository.SocialGraphRepository
	contentRepo        repository.ContentRepository
	userRPCClient      UserRPCClient          // 调用 user 服务获取用户信息
	reasonConfigClient ReasonTextConfigClient // 调用配置服务获取推荐理由文案（可选）
}

// UserRPCClient 用户服务RPC客户端接口
// 定义在应用层，因为这是技术细节
type UserRPCClient interface {
	GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error)
	GetUserInfoBatch(ctx context.Context, userIDs []int64) ([]*UserInfo, error)
}

// ReasonTextConfigClient 推荐理由文案配置服务客户端接口
// 用于从配置服务获取推荐理由的展示文案
type ReasonTextConfigClient interface {
	// GetReasonText 获取推荐理由的展示文案
	// reasonType: 推荐理由类型（如 "followed_by_following"）
	// count: 相关用户数量（用于生成文案，如 "3 位你关注的人"）
	// 返回配置的文案，如果配置服务异常或没有配置，返回空字符串（会降级到本地逻辑）
	GetReasonText(ctx context.Context, reasonType string, count int) (string, error)
}

// UserInfo 用户信息（来自 user 服务）
type UserInfo struct {
	UserID   int64
	Username string
	Avatar   string
	Bio      string
}

// NewRecommendationService 构造函数
// reasonConfigClient 可以为 nil，表示不使用配置服务（降级到本地逻辑）
func NewRecommendationService(
	generator *service.RecommendationGenerator,
	socialGraphRepo repository.SocialGraphRepository,
	contentRepo repository.ContentRepository,
	userRPCClient UserRPCClient,
	reasonConfigClient ReasonTextConfigClient,
) *RecommendationService {
	return &RecommendationService{
		generator:          generator,
		socialGraphRepo:    socialGraphRepo,
		contentRepo:        contentRepo,
		userRPCClient:      userRPCClient,
		reasonConfigClient: reasonConfigClient,
	}
}

// GetFollowingBasedRecommendations 用例：获取基于关注的推荐
//
// 这是一个完整的业务用例（Use Case），展示了应用服务如何编排。
//
// 用例流程：
// 1. 参数转换：int64 → 领域对象（UserID）
// 2. 调用领域服务：生成推荐列表
// 3. 获取 Top N：按分数排序取前 N 个
// 4. 批量获取用户信息：调用 user 服务（性能优化）
// 5. 获取用户帖子：调用 content 服务
// 6. 组装响应：领域对象 → DTO
//
// 为什么这些逻辑在应用层？
// - 跨服务调用：涉及技术细节（RPC）
// - 性能优化：批量查询是技术决策
// - DTO 转换：适配外部接口
// 这些都不是核心业务规则，所以不在领域层。
//
// 实际业务场景：
// 用户打开"推荐关注"页面 →
//
//	前端调用这个接口 →
//	返回推荐用户列表（包含头像、简介、最近帖子）
//
// 性能考虑：
// - 批量获取用户信息：避免 N+1 查询问题
// - 容错处理：某个用户信息获取失败不影响整体
// - 限制数量：通过 limit 参数控制返回数量
func (s *RecommendationService) GetFollowingBasedRecommendations(
	ctx context.Context,
	userID int64,
	limit int,
) (*dto.RecommendationResponse, error) {

	// 步骤1：转换为领域对象
	domainUserID, err := valueobject.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	// 步骤2：调用领域服务生成推荐
	recommendationList, err := s.generator.GenerateFollowingBasedRecommendations(
		ctx, domainUserID, 7, // 最近7天
	)
	if err != nil {
		return nil, err
	}

	// 步骤3：获取 Top N 推荐
	topRecommendations := recommendationList.GetTopN(limit)

	// 如果没有推荐，直接返回空列表
	if len(topRecommendations) == 0 {
		return &dto.RecommendationResponse{
			Recommendations: []*dto.UserRecommendationDTO{},
		}, nil
	}

	// 步骤4：批量获取用户信息（优化性能）
	userIDs := make([]int64, 0, len(topRecommendations))
	for _, rec := range topRecommendations {
		userIDs = append(userIDs, rec.TargetUserID().Value())
	}

	userInfoMap, err := s.getUserInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	// 步骤5：组装响应数据
	response := &dto.RecommendationResponse{
		Recommendations: make([]*dto.UserRecommendationDTO, 0, len(topRecommendations)),
	}

	for _, rec := range topRecommendations {
		// 获取用户详情
		userInfo, exists := userInfoMap[rec.TargetUserID().Value()]
		if !exists {
			continue // 跳过无法获取信息的用户
		}

		// 获取用户最近的帖子
		posts, err := s.contentRepo.GetRecentPosts(ctx, rec.TargetUserID(), 3)
		if err != nil {
			posts = nil // 容错：获取失败不影响推荐
		}

		// 获取推荐理由文案（优先使用配置服务）
		reasonText := s.getReasonText(ctx, rec.Reason())

		// 转换为 DTO
		recommendationDTO := &dto.UserRecommendationDTO{
			UserID:      rec.TargetUserID().Value(),
			Username:    userInfo.Username,
			Avatar:      userInfo.Avatar,
			Bio:         userInfo.Bio,
			Reason:      reasonText,
			Score:       rec.Score(),
			RecentPosts: s.convertPostsToDTO(posts),
		}

		response.Recommendations = append(response.Recommendations, recommendationDTO)
	}

	return response, nil
}

// getUserInfoMap 辅助方法：批量获取用户信息并转换为 map
func (s *RecommendationService) getUserInfoMap(
	ctx context.Context,
	userIDs []int64,
) (map[int64]*UserInfo, error) {
	userInfos, err := s.userRPCClient.GetUserInfoBatch(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]*UserInfo, len(userInfos))
	for _, info := range userInfos {
		result[info.UserID] = info
	}
	return result, nil
}

// convertPostsToDTO 辅助方法：转换帖子实体为 DTO
func (s *RecommendationService) convertPostsToDTO(posts []*entity.Post) []*dto.PostDTO {
	if posts == nil {
		return []*dto.PostDTO{}
	}

	result := make([]*dto.PostDTO, 0, len(posts))
	for _, post := range posts {
		result = append(result, &dto.PostDTO{
			PostID:    post.ID().Value(),
			Content:   post.Content(),
			CreatedAt: post.CreatedAt().Format("2006-01-02 15:04:05"),
		})
	}
	return result
}

// getReasonText 辅助方法：获取推荐理由文案
//
// 这个方法展示了如何在应用层集成配置服务，同时保持降级能力。
//
// 设计思路：
// 1. 优先尝试从配置服务获取文案（如果配置了 reasonConfigClient）
// 2. 如果配置服务不可用或返回空，降级到领域对象的本地逻辑
// 3. 保证无论配置服务是否可用，都能正常展示推荐理由
//
// 为什么在应用层而不是领域层？
// - 调用 HTTP 服务是基础设施细节，不应该污染领域层
// - 配置文案更像是展示层的关注点，不是核心业务规则
// - 应用层负责协调外部服务，这正是它的职责
//
// 实际场景：
//
//	// 场景1：配置服务正常
//	reasonConfigClient 返回 "你的 3 位好友也关注了TA"
//	→ 使用配置服务的文案
//
//	// 场景2：配置服务异常或未配置
//	reasonConfigClient 为 nil 或返回错误
//	→ 降级到 reason.Description()（本地逻辑）
//
//	// 场景3：配置服务返回空字符串
//	reasonConfigClient 返回 ""
//	→ 降级到 reason.Description()（本地逻辑）
//
// 容错设计：
// - reasonConfigClient 可以为 nil（表示不使用配置服务）
// - 配置服务调用失败不影响推荐功能
// - 配置服务返回空字符串时降级到本地逻辑
//
// 扩展性：
// 未来可以添加更多逻辑：
// - 缓存配置文案（减少 HTTP 调用）
// - A/B 测试（根据用户分组返回不同文案）
// - 多语言支持（根据用户语言返回对应文案）
func (s *RecommendationService) getReasonText(ctx context.Context, reason valueobject.RecommendationReason) string {
	// 如果没有配置客户端，直接使用本地逻辑
	if s.reasonConfigClient == nil {
		return reason.Description()
	}

	// 将领域对象的类型转换为配置服务的类型标识
	var reasonType string
	switch reason.Type() {
	case valueobject.ReasonFollowedByFollowing:
		reasonType = "followed_by_following"
	case valueobject.ReasonPopularInNetwork:
		reasonType = "popular_in_network"
	default:
		reasonType = "default"
	}

	// 尝试从配置服务获取文案
	configText, err := s.reasonConfigClient.GetReasonText(
		ctx,
		reasonType,
		len(reason.RelatedUsers()),
	)

	// 容错处理：配置服务异常或返回空，降级到本地逻辑
	if err != nil || configText == "" {
		return reason.Description()
	}

	return configText
}
