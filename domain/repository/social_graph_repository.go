package repository

import (
	"context"

	"service/domain/valueobject"
)

// SocialGraphRepository 仓储接口：社交关系图
//
// 什么是仓储模式？
// 仓储（Repository）是 DDD 中用于访问聚合的模式，
// 它提供类似集合的接口来访问和持久化聚合。
//
// 为什么接口定义在领域层？（依赖倒置原则 DIP）
//
// 传统分层架构：
//
//	领域层 → 依赖 → 数据访问层
//	问题：领域层依赖技术细节，难以测试和替换
//
// DDD 架构（依赖倒置）：
//
//	领域层：定义仓储接口（我需要什么数据）
//	   ↑
//	基础设施层：实现仓储接口（如何获取数据）
//	好处：领域层不依赖任何外层，保持纯粹
//
// 实际好处：
// 1. 可测试性：单元测试时可以用 mock 实现替换真实数据库
// 2. 可替换性：从 MySQL 切换到 MongoDB 不影响领域层
// 3. 业务语言：接口方法名反映业务概念，不是技术术语
//
// 仓储 vs DAO：
// - DAO：面向数据表，方法如 findByFollowerId
// - Repository：面向聚合，方法如 GetFollowings（业务语言）
//
// 使用示例：
//
//	// 领域层代码
//	followings, err := repo.GetFollowings(ctx, userID)
//	// 不关心数据来自 MySQL 还是 Redis
type SocialGraphRepository interface {
	// GetFollowings 获取用户关注的所有人
	//
	// 业务含义：查询用户的关注列表
	// 返回：用户ID列表
	GetFollowings(ctx context.Context, userID valueobject.UserID) ([]valueobject.UserID, error)

	// GetRecentFollowings 获取用户最近N天关注的人
	//
	// 业务含义：查询用户最近的关注行为
	// 参数：
	// - userID: 用户ID
	// - days: 最近多少天
	// 返回：用户ID列表
	GetRecentFollowings(ctx context.Context, userID valueobject.UserID, days int) ([]valueobject.UserID, error)

	// IsFollowing 检查用户A是否关注了用户B
	//
	// 业务含义：判断关注关系是否存在
	IsFollowing(ctx context.Context, followerID, followingID valueobject.UserID) (bool, error)
}
