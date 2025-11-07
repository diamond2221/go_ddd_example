package dto

// DTO（数据传输对象 - Data Transfer Object）
//
// 什么是 DTO？
// DTO 是用于在不同层之间传输数据的简单对象，
// 通常只包含数据字段，没有业务逻辑。
//
// 为什么需要 DTO？为什么不直接返回领域对象？
//
// 问题场景：
// 如果直接返回领域对象 UserRecommendation：
// 1. 暴露内部实现：客户端看到所有内部字段
// 2. 难以序列化：领域对象可能包含复杂的关联
// 3. 版本耦合：领域模型改变会破坏 API 契约
// 4. 性能问题：可能传输不需要的数据
//
// DTO 的好处：
// 1. 解耦：领域对象变化不影响 API 契约
//    - 领域层重构不影响客户端
//    - 可以独立演进内部模型
//
// 2. 简化：只传输需要的数据
//    - 减少网络传输量
//    - 客户端更容易理解
//
// 3. 适配：适配不同的展示需求
//    - Web 端可能需要更多字段
//    - Mobile 端可能需要更少字段
//    - 可以为不同客户端定义不同的 DTO
//
// 4. 安全：不暴露敏感信息
//    - 可以过滤掉内部字段
//    - 控制数据访问权限
//
// DTO vs 领域对象 vs PO：
// ┌──────────┬────────────┬──────────────┬──────────────┐
// │          │ 领域对象    │ DTO          │ PO           │
// ├──────────┼────────────┼──────────────┼──────────────┤
// │ 位置     │ 领域层      │ 应用层/接口层 │ 基础设施层    │
// │ 职责     │ 业务逻辑    │ 数据传输      │ 数据持久化    │
// │ 行为     │ 有          │ 无            │ 无            │
// │ 依赖     │ 无外部依赖  │ 无外部依赖    │ 依赖 ORM     │
// │ 示例     │ UserRec...  │ UserRec...DTO │ FollowPO     │
// └──────────┴────────────┴──────────────┴──────────────┘
//
// 实际使用流程：
// 数据库（PO）→ 仓储转换 → 领域对象 → 应用服务转换 → DTO → JSON → 客户端
//
// 代价：
// - 需要写转换代码
// - 代码量增加
//
// 何时值得使用 DTO？
// - 对外 API：必须使用，保护内部实现
// - 内部服务：可以考虑直接用领域对象（如果信任内部调用）

// RecommendationResponse 推荐响应
type RecommendationResponse struct {
	Recommendations []*UserRecommendationDTO `json:"recommendations"`
}

// UserRecommendationDTO 用户推荐DTO
type UserRecommendationDTO struct {
	UserID      int64      `json:"user_id"`
	Username    string     `json:"username"`
	Avatar      string     `json:"avatar"`
	Bio         string     `json:"bio"`
	Reason      string     `json:"reason"`       // "3 位你关注的人也关注了TA"
	Score       int        `json:"score"`        // 推荐分数
	RecentPosts []*PostDTO `json:"recent_posts"` // 最近的帖子
}

// PostDTO 帖子DTO
type PostDTO struct {
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"` // 格式化后的时间字符串
}
