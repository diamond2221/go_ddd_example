package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"service/domain/repository"
	"service/domain/valueobject"
)

// SocialGraphRepositoryImpl 仓储实现（基础设施层）
//
// 这是仓储模式的实现端，展示了依赖倒置原则的实践。
//
// 职责：
// 1. 实现领域层定义的仓储接口
// 2. 处理数据库访问细节（SQL、事务、连接池等）
// 3. 领域对象 ↔ 持久化对象（PO）转换
// 4. 性能优化（索引、缓存、批量查询等）
//
// 依赖方向：
//
//	领域层：定义 SocialGraphRepository 接口
//	   ↑
//	基础设施层：实现 SocialGraphRepositoryImpl
//
// 为什么这样设计？
// 传统方式：领域层依赖数据访问层
// DDD 方式：数据访问层依赖领域层（依赖倒置）
//
// 好处：
// 1. 领域层独立：不依赖具体的数据库技术
// 2. 可测试性：可以用 mock 实现替换真实数据库
// 3. 可替换性：从 MySQL 切换到 MongoDB 不影响领域层
//
// 实际场景：
// 领域服务调用：repo.GetFollowings(ctx, userID)
// 仓储实现：
//  1. 构造 SQL 查询
//  2. 执行数据库查询
//  3. 将 FollowPO 转换为 UserID
//  4. 返回给领域层
type SocialGraphRepositoryImpl struct {
	db *gorm.DB
}

// NewSocialGraphRepository 构造函数
// 返回接口类型，而不是具体类型
func NewSocialGraphRepository(db *gorm.DB) repository.SocialGraphRepository {
	return &SocialGraphRepositoryImpl{db: db}
}

// GetFollowings 实现接口：获取用户关注的所有人
//
// 这个方法展示了仓储实现的典型模式：
// 1. 使用 ORM 查询数据库
// 2. 将持久化对象（PO）转换为领域对象
// 3. 返回领域对象给调用者
//
// 注意事项：
// - 使用 ctx 支持超时和取消
// - 只查询 status = 'active' 的关注关系（软删除）
// - 转换时忽略错误（实际项目中应该记录日志）
//
// 性能优化点：
// - 可以添加缓存（Redis）
// - 可以添加索引（idx_follower）
// - 可以分页查询（如果关注数很多）
func (r *SocialGraphRepositoryImpl) GetFollowings(
	ctx context.Context,
	userID valueobject.UserID,
) ([]valueobject.UserID, error) {

	var follows []FollowPO
	err := r.db.WithContext(ctx).
		Where("follower_id = ? AND status = ?", userID.Value(), "active").
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	// 转换 PO -> 领域对象
	// 这是仓储的重要职责：隔离数据库模型和领域模型
	result := make([]valueobject.UserID, 0, len(follows))
	for _, follow := range follows {
		domainID, _ := valueobject.NewUserID(follow.FollowingID)
		result = append(result, domainID)
	}

	return result, nil
}

// GetRecentFollowings 实现接口：获取用户最近N天关注的人
func (r *SocialGraphRepositoryImpl) GetRecentFollowings(
	ctx context.Context,
	userID valueobject.UserID,
	days int,
) ([]valueobject.UserID, error) {

	since := time.Now().AddDate(0, 0, -days)

	var follows []FollowPO
	err := r.db.WithContext(ctx).
		Where("follower_id = ? AND status = ? AND created_at >= ?",
			userID.Value(), "active", since).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	// 转换 PO -> 领域对象
	result := make([]valueobject.UserID, 0, len(follows))
	for _, follow := range follows {
		domainID, _ := valueobject.NewUserID(follow.FollowingID)
		result = append(result, domainID)
	}

	return result, nil
}

// IsFollowing 实现接口：检查关注关系
func (r *SocialGraphRepositoryImpl) IsFollowing(
	ctx context.Context,
	followerID, followingID valueobject.UserID,
) (bool, error) {

	var count int64
	err := r.db.WithContext(ctx).
		Model(&FollowPO{}).
		Where("follower_id = ? AND following_id = ? AND status = ?",
			followerID.Value(), followingID.Value(), "active").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FollowPO 持久化对象（PO - Persistent Object）
//
// 为什么需要 PO？为什么不直接用领域对象？
//
// 这是 DDD 中的一个重要设计决策：领域对象与持久化对象分离
//
// 领域对象的特点：
// - 表达业务概念
// - 包含业务行为
// - 不依赖技术框架
// - 例如：UserRecommendation 有 CalculateScore() 方法
//
// 持久化对象的特点：
// - 映射数据库表结构
// - 包含 ORM 标签（gorm、json 等）
// - 只有数据，没有行为
// - 例如：FollowPO 对应 follows 表
//
// 分离的好处：
// 1. 领域模型独立：不受数据库结构影响
//   - 数据库表改名不影响领域代码
//   - 可以从 GORM 切换到其他 ORM
//
// 2. 灵活的映射：
//   - 一个领域对象可能对应多个表
//   - 多个领域对象可能来自一个表
//
// 3. 性能优化：
//   - 查询时可以只加载需要的字段
//   - 不需要加载领域对象的所有关联数据
//
// 4. 测试友好：
//   - 领域对象可以纯内存测试
//   - 不需要数据库就能测试业务逻辑
//
// 代价：
// - 需要写转换代码（PO ↔ 领域对象）
// - 代码量增加
//
// 何时值得分离？
// - 复杂业务系统：值得
// - 简单 CRUD 应用：可能过度设计
//
// 实际场景：
// 数据库的 follows 表可能有很多字段（created_by, updated_by 等），
// 但领域层只关心核心的关注关系，不需要这些技术字段。
type FollowPO struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	FollowerID  int64     `gorm:"index:idx_follower;not null"`
	FollowingID int64     `gorm:"index:idx_following;not null"`
	Status      string    `gorm:"type:varchar(20);default:'active'"`
	CreatedAt   time.Time `gorm:"index:idx_created_at;not null"`
	UpdatedAt   time.Time
}

// TableName 指定表名
func (FollowPO) TableName() string {
	return "follows"
}
