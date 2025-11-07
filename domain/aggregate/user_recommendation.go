package aggregate

import (
	"errors"
	"time"

	"service/domain/valueobject"
)

var (
	ErrNoReasonForRecommendation = errors.New("no reason for recommendation")
)

// UserRecommendation 聚合根：用户推荐
//
// 什么是聚合？
// 聚合是一组相关对象的集合，作为一个整体来管理。
// 聚合根是聚合的入口，外部只能通过聚合根来访问聚合内的对象。
//
// 为什么 UserRecommendation 是聚合根？
// 1. 完整的生命周期：推荐有创建、使用、过期的完整生命周期
// 2. 业务规则封装：推荐分数计算、过期判断等规则都在这里
// 3. 一致性边界：推荐的所有数据必须保持一致（分数、理由、时间等）
// 4. 事务边界：对推荐的修改应该在一个事务内完成
//
// 聚合 vs 实体 vs 值对象：
// - 值对象（UserID）：无标识，不可变，如 123
// - 实体（Post）：有标识，可变，如某个具体的帖子
// - 聚合根（UserRecommendation）：有标识，管理一组对象，如一个完整的推荐
//
// 设计原则：
// 1. 聚合应该尽可能小：只包含必须保持一致性的数据
// 2. 通过聚合根修改：外部不能直接修改聚合内部的对象
// 3. 聚合之间通过 ID 引用：不直接持有其他聚合的引用
//
// 实际业务场景：
// 当用户刷新推荐列表时，系统会创建多个 UserRecommendation 对象，
// 每个对象代表一个推荐用户，包含推荐理由、分数等完整信息。
type UserRecommendation struct {
	// 私有字段，只能通过方法访问，保证封装性
	id              valueobject.RecommendationID
	targetUserID    valueobject.UserID // 被推荐的用户
	reason          valueobject.RecommendationReason
	score           int       // 推荐分数
	recentPostCount int       // 最近帖子数
	createdAt       time.Time // 创建时间
	expiresAt       time.Time // 过期时间
}

// NewUserRecommendation 工厂方法：创建新的用户推荐
//
// 为什么用工厂方法而不是直接 new？
// 工厂方法是 DDD 中创建聚合的标准方式，好处：
// 1. 集中验证：所有创建时的业务规则都在这里
// 2. 保证有效性：创建成功的对象一定是有效的
// 3. 封装复杂性：隐藏创建逻辑的复杂性
//
// 在创建时执行的业务规则：
// 1. 必须有推荐理由（至少1个关注者）
// 2. 自动计算推荐分数（根据关注者数和帖子数）
// 3. 设置过期时间（7天后过期）
// 4. 生成唯一的推荐ID
//
// 使用示例：
//
//	reason := valueobject.NewFollowedByFollowingReason([]UserID{user1, user2})
//	rec, err := NewUserRecommendation(targetUser, reason, 5)
//	if err != nil {
//	    // 处理创建失败（如没有推荐理由）
//	}
//	// rec 保证是有效的推荐对象
//
// 对比直接 new：
//
//	rec := &UserRecommendation{...} // 可能忘记验证，可能忘记计算分数
//	工厂方法保证了对象的完整性和有效性
func NewUserRecommendation(
	targetUserID valueobject.UserID,
	reason valueobject.RecommendationReason,
	recentPostCount int,
) (*UserRecommendation, error) {
	// 业务规则：至少要有1个关注者才能推荐
	if len(reason.RelatedUsers()) == 0 {
		return nil, ErrNoReasonForRecommendation
	}

	// 业务规则：计算推荐分数
	score := calculateScore(reason, recentPostCount)

	now := time.Now()
	return &UserRecommendation{
		id:              valueobject.NewRecommendationID(),
		targetUserID:    targetUserID,
		reason:          reason,
		score:           score,
		recentPostCount: recentPostCount,
		createdAt:       now,
		expiresAt:       now.Add(7 * 24 * time.Hour), // 7天过期
	}, nil
}

// calculateScore 业务规则：推荐分数计算
//
// 这是核心业务规则，决定了推荐的排序。
//
// 计算公式：
// - 基础分数 = 推荐理由权重（关注者数 × 10）
// - 活跃度加分 = 帖子数量 × 2
//
// 业务逻辑：
// - 被更多人关注的用户分数更高
// - 有活跃内容的用户更值得推荐
//
// 实际示例：
//
//	用户A：3个关注者，5个帖子 → 分数 = 3×10 + 5×2 = 40
//	用户B：1个关注者，10个帖子 → 分数 = 1×10 + 10×2 = 30
//	结果：优先推荐用户A（社交信号更强）
//
// 为什么这个逻辑在领域层？
// 因为这是核心业务规则，产品经理定义的推荐策略。
// 如果策略改变（如调整权重），只需修改这里。
//
// 扩展性：
// 未来可以添加更多因素：
// - 用户活跃度（最后登录时间）
// - 内容质量（点赞数、评论数）
// - 个性化因素（兴趣匹配度）
func calculateScore(reason valueobject.RecommendationReason, postCount int) int {
	score := reason.Weight()

	// 有活跃内容加分
	if postCount > 0 {
		score += postCount * 2
	}

	return score
}

// IsExpired 业务规则：推荐是否过期
//
// 过期策略：
// - 推荐生成后 7 天过期
// - 过期的推荐不应该再展示给用户
func (r *UserRecommendation) IsExpired() bool {
	return time.Now().After(r.expiresAt)
}

// --- 访问器方法（Getters）---
// 聚合根对外暴露的只读接口
// 保证外部无法直接修改内部状态

func (r *UserRecommendation) ID() valueobject.RecommendationID {
	return r.id
}

func (r *UserRecommendation) TargetUserID() valueobject.UserID {
	return r.targetUserID
}

func (r *UserRecommendation) Reason() valueobject.RecommendationReason {
	return r.reason
}

func (r *UserRecommendation) Score() int {
	return r.score
}

func (r *UserRecommendation) RecentPostCount() int {
	return r.recentPostCount
}

func (r *UserRecommendation) CreatedAt() time.Time {
	return r.createdAt
}

func (r *UserRecommendation) ExpiresAt() time.Time {
	return r.expiresAt
}

// --- 领域行为方法 ---
// 如果需要修改推荐，应该通过这些方法
// 而不是直接修改字段

// Refresh 业务行为：刷新推荐（延长过期时间）
func (r *UserRecommendation) Refresh() {
	r.expiresAt = time.Now().Add(7 * 24 * time.Hour)
}

// UpdatePostCount 业务行为：更新帖子数量并重新计算分数
func (r *UserRecommendation) UpdatePostCount(newCount int) {
	r.recentPostCount = newCount
	r.score = calculateScore(r.reason, newCount)
}
