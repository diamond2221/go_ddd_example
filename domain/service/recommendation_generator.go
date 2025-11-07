package service

import (
	"context"

	"service/domain/repository"

	"service/domain/aggregate"
	"service/domain/valueobject"
)

// RecommendationGenerator 领域服务：推荐生成逻辑
//
// 什么是领域服务？
// 领域服务是领域层的一种模式，用于处理不属于任何单一实体或聚合的业务逻辑。
//
// 什么时候使用领域服务？
// 1. 业务逻辑涉及多个聚合
// 2. 业务逻辑不自然地属于某个实体
// 3. 业务逻辑是核心领域知识
//
// 为什么推荐生成是领域服务？
// 推荐生成需要：
// - 用户的关注关系（SocialGraph 聚合）
// - 用户的内容数据（Content 聚合）
// - 推荐算法逻辑（核心业务规则）
// 这些逻辑不属于任何单一聚合，所以放在领域服务中。
//
// 领域服务 vs 应用服务：
// ┌─────────────────┬──────────────────┬──────────────────┐
// │                 │   领域服务        │   应用服务        │
// ├─────────────────┼──────────────────┼──────────────────┤
// │ 位置            │   领域层          │   应用层          │
// │ 职责            │   纯业务逻辑      │   用例编排        │
// │ 依赖            │   仓储接口        │   领域服务+RPC    │
// │ 示例            │   推荐算法        │   获取推荐用例    │
// └─────────────────┴──────────────────┴──────────────────┘
//
// 实际业务场景：
// 产品经理说："我们要推荐用户关注的人最近关注的人"
// 这个算法逻辑就是领域服务要实现的核心业务规则。
//
// 对比传统方式：
// 传统方式：这些逻辑可能散落在 Service 层的各个方法中
// DDD 方式：集中在领域服务中，清晰表达业务意图
type RecommendationGenerator struct {
	socialGraphRepo repository.SocialGraphRepository
	contentRepo     repository.ContentRepository
}

// NewRecommendationGenerator 构造函数
func NewRecommendationGenerator(
	socialGraphRepo repository.SocialGraphRepository,
	contentRepo repository.ContentRepository,
) *RecommendationGenerator {
	return &RecommendationGenerator{
		socialGraphRepo: socialGraphRepo,
		contentRepo:     contentRepo,
	}
}

// GenerateFollowingBasedRecommendations 核心领域逻辑：生成基于关注的推荐
//
// 这是推荐算法的核心实现，体现了业务规则。
//
// 业务需求（产品经理的话）：
// "推荐我关注的人最近关注的用户，按关注者数量和活跃度排序"
//
// 算法流程：
// 1. 获取用户A关注的人（B、C、D）
// 2. 获取B、C、D最近N天关注的人（E、F、G）
// 3. 统计每个被关注用户的关注者数量
// 4. 获取被关注用户的活跃度（帖子数）
// 5. 计算推荐分数并创建推荐对象
// 6. 返回推荐列表（会自动去重、过滤自己）
//
// 实际示例：
//
//	用户A关注了 [B, C, D]
//	B最近关注了 [E, F]
//	C最近关注了 [E, G]
//	D最近关注了 [H]
//	结果：推荐 E（2人关注）、F（1人）、G（1人）、H（1人）
//
// 为什么这个逻辑在领域服务？
// 1. 涉及多个聚合（用户、关注关系、内容）
// 2. 是核心业务规则，不是技术细节
// 3. 不属于任何单一聚合的职责
//
// 容错设计：
// - 某个用户数据获取失败不影响整体
// - 帖子数获取失败默认为0
// - 无效推荐会被跳过
//
// 参数：
// - forUserID: 为哪个用户生成推荐
// - days: 最近多少天的关注（通常是7天）
func (g *RecommendationGenerator) GenerateFollowingBasedRecommendations(
	ctx context.Context,
	forUserID valueobject.UserID,
	days int,
) (*aggregate.RecommendationList, error) {

	// 创建推荐列表聚合
	list := aggregate.NewRecommendationList(forUserID)

	// 步骤1：获取用户关注的人
	followings, err := g.socialGraphRepo.GetFollowings(ctx, forUserID)
	if err != nil {
		return nil, err
	}

	// 如果用户没有关注任何人，返回空列表
	if len(followings) == 0 {
		return list, nil
	}

	// 步骤2：获取这些人最近关注的人（去重）
	// key: 被关注的用户ID
	// value: 哪些用户关注了这个人
	recentFollowedUsers := make(map[valueobject.UserID][]valueobject.UserID)

	for _, following := range followings {
		// 获取这个用户最近关注的人
		recentFollows, err := g.socialGraphRepo.GetRecentFollowings(
			ctx, following, days,
		)
		if err != nil {
			// 容错处理：某个用户的数据获取失败不影响整体
			continue
		}

		// 记录谁关注了谁
		for _, newFollow := range recentFollows {
			recentFollowedUsers[newFollow] = append(
				recentFollowedUsers[newFollow],
				following,
			)
		}
	}

	// 步骤3：为每个推荐用户创建推荐对象
	for targetUserID, followedBy := range recentFollowedUsers {
		// 获取该用户最近的帖子数
		postCount, err := g.contentRepo.CountRecentPosts(ctx, targetUserID, days)
		if err != nil {
			postCount = 0 // 容错：获取失败默认为0
		}

		// 创建推荐理由
		reason := valueobject.NewFollowedByFollowingReason(followedBy)

		// 创建推荐聚合
		recommendation, err := aggregate.NewUserRecommendation(
			targetUserID,
			reason,
			postCount,
		)
		if err != nil {
			// 跳过无效推荐（如没有推荐理由）
			continue
		}

		// 添加到推荐列表
		if err := list.AddRecommendation(recommendation); err != nil {
			// 跳过重复或无效推荐（如推荐自己）
			continue
		}
	}

	return list, nil
}

// GeneratePopularityBasedRecommendations 扩展示例：基于热度的推荐
//
// 这展示了如何扩展新的推荐策略：
// 1. 在同一个领域服务中添加新方法
// 2. 或者创建新的领域服务类
func (g *RecommendationGenerator) GeneratePopularityBasedRecommendations(
	ctx context.Context,
	forUserID valueobject.UserID,
) (*aggregate.RecommendationList, error) {
	// TODO: 实现基于热度的推荐逻辑
	// 例如：推荐在用户社交网络中被多人关注的用户
	return aggregate.NewRecommendationList(forUserID), nil
}
