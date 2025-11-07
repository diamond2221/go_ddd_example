package aggregate

import (
	"errors"
	"sort"
	"time"

	"service/domain/valueobject"
)

var (
	ErrCannotRecommendSelf     = errors.New("cannot recommend self")
	ErrDuplicateRecommendation = errors.New("duplicate recommendation")
)

// RecommendationList 聚合：推荐列表
//
// 为什么需要单独的推荐列表聚合？
// 这是一个常见的设计问题：是否需要为集合创建单独的聚合？
//
// 答案取决于：集合本身是否有业务规则？
//
// 在这个场景中，推荐列表有自己的业务规则：
// 1. 去重规则：不能推荐同一个用户两次
// 2. 排序规则：按分数排序
// 3. 过滤规则：移除过期推荐、低分推荐
// 4. 约束规则：不能推荐自己
//
// 如果只是简单的 []UserRecommendation，这些规则会散落在各处。
// 通过创建 RecommendationList 聚合，我们把这些规则集中管理。
//
// 与 UserRecommendation 的关系：
// - 它们是独立的聚合，不是父子关系
// - RecommendationList 持有 UserRecommendation 的引用（指针）
// - 每个聚合维护自己的一致性边界
//
// 实际业务场景：
// 用户打开推荐页面时，系统生成一个 RecommendationList，
// 包含多个 UserRecommendation。列表会自动去重、排序、过滤，
// 保证用户看到的是高质量的推荐。
//
// 对比传统方式：
// 传统方式：在 Service 层用循环和 if 判断处理这些逻辑
// DDD 方式：在聚合中封装这些业务规则，代码更清晰
type RecommendationList struct {
	forUserID       valueobject.UserID    // 为哪个用户生成的推荐
	recommendations []*UserRecommendation // 推荐列表
	generatedAt     time.Time             // 生成时间
}

// NewRecommendationList 工厂方法：创建新的推荐列表
func NewRecommendationList(forUserID valueobject.UserID) *RecommendationList {
	return &RecommendationList{
		forUserID:       forUserID,
		recommendations: make([]*UserRecommendation, 0),
		generatedAt:     time.Now(),
	}
}

// AddRecommendation 业务行为：添加推荐
//
// 这个方法展示了聚合如何保护业务不变量（Invariants）。
//
// 业务不变量：
// 1. 不能推荐自己（产品规则：自己不需要关注自己）
// 2. 不能重复推荐（产品规则：同一用户只推荐一次）
//
// 为什么在聚合中验证？
// 如果在外部验证，可能会遗漏或不一致。
// 聚合保证：无论谁调用，这些规则都会被执行。
//
// 实际场景：
//
//	list := NewRecommendationList(userA)
//	list.AddRecommendation(recA) // 成功
//	list.AddRecommendation(recA) // 失败：重复推荐
//	list.AddRecommendation(recSelf) // 失败：推荐自己
//
// 对比传统方式：
// 传统方式：在 Service 层用 if 判断，容易遗漏
// DDD 方式：在聚合中强制执行，保证一致性
func (l *RecommendationList) AddRecommendation(rec *UserRecommendation) error {
	// 业务规则：不能推荐自己
	if rec.TargetUserID().Equals(l.forUserID) {
		return ErrCannotRecommendSelf
	}

	// 业务规则：不能重复推荐
	for _, existing := range l.recommendations {
		if existing.TargetUserID().Equals(rec.TargetUserID()) {
			return ErrDuplicateRecommendation
		}
	}

	l.recommendations = append(l.recommendations, rec)
	return nil
}

// GetTopN 业务行为：获取分数最高的 N 个推荐
//
// 这是一个查询方法，展示了聚合如何封装业务逻辑。
//
// 业务规则：
// 1. 按分数降序排序（分数高的优先展示）
// 2. 返回前 N 个（控制展示数量）
// 3. 如果总数不足 N，返回全部
//
// 为什么在聚合中排序？
// 排序规则是业务规则的一部分，应该由聚合控制。
// 外部调用者不需要知道如何排序，只需要"给我最好的N个"。
//
// 实际场景：
//
//	list.AddRecommendation(rec1) // 分数 40
//	list.AddRecommendation(rec2) // 分数 30
//	list.AddRecommendation(rec3) // 分数 50
//	top2 := list.GetTopN(2) // 返回 [rec3(50), rec1(40)]
//
// 设计考虑：
// - 返回副本：不修改原列表，避免副作用
// - 性能：每次调用都排序，如果频繁调用可以优化（缓存排序结果）
func (l *RecommendationList) GetTopN(n int) []*UserRecommendation {
	// 创建副本进行排序，不修改原列表
	sorted := make([]*UserRecommendation, len(l.recommendations))
	copy(sorted, l.recommendations)

	// 按分数降序排序
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score() > sorted[j].Score()
	})

	// 返回前 N 个
	if len(sorted) > n {
		return sorted[:n]
	}
	return sorted
}

// RemoveExpired 业务行为：移除过期推荐
//
// 业务规则：
// - 过期的推荐不应该再展示给用户
// - 定期清理过期推荐，保持列表干净
func (l *RecommendationList) RemoveExpired() {
	valid := make([]*UserRecommendation, 0)
	for _, rec := range l.recommendations {
		if !rec.IsExpired() {
			valid = append(valid, rec)
		}
	}
	l.recommendations = valid
}

// FilterByMinScore 业务行为：过滤低分推荐
//
// 业务规则：
// - 只保留分数达到最低标准的推荐
// - 提高推荐质量
func (l *RecommendationList) FilterByMinScore(minScore int) {
	filtered := make([]*UserRecommendation, 0)
	for _, rec := range l.recommendations {
		if rec.Score() >= minScore {
			filtered = append(filtered, rec)
		}
	}
	l.recommendations = filtered
}

// Count 查询方法：获取推荐数量
func (l *RecommendationList) Count() int {
	return len(l.recommendations)
}

// IsEmpty 查询方法：列表是否为空
func (l *RecommendationList) IsEmpty() bool {
	return len(l.recommendations) == 0
}

// ForUserID 访问器：获取目标用户ID
func (l *RecommendationList) ForUserID() valueobject.UserID {
	return l.forUserID
}

// GeneratedAt 访问器：获取生成时间
func (l *RecommendationList) GeneratedAt() time.Time {
	return l.generatedAt
}

// All 访问器：获取所有推荐（返回副本）
func (l *RecommendationList) All() []*UserRecommendation {
	result := make([]*UserRecommendation, len(l.recommendations))
	copy(result, l.recommendations)
	return result
}
