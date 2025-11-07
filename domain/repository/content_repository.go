package repository

import (
	"context"

	"service/domain/entity"
	"service/domain/valueobject"
)

// ContentRepository 仓储接口：内容数据
//
// 注意：这里的 Post 是领域实体，不是数据库模型
type ContentRepository interface {
	// CountRecentPosts 统计用户最近N天的帖子数
	//
	// 业务含义：评估用户的活跃度
	// 用于推荐分数计算
	CountRecentPosts(ctx context.Context, userID valueobject.UserID, days int) (int, error)

	// GetRecentPosts 获取用户最近的帖子
	//
	// 业务含义：展示推荐用户的内容
	// 参数：
	// - userID: 用户ID
	// - limit: 最多返回多少条
	GetRecentPosts(ctx context.Context, userID valueobject.UserID, limit int) ([]*entity.Post, error)
}
