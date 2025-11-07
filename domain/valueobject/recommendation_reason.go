package valueobject

import "fmt"

// ReasonType 推荐理由类型
type ReasonType int

const (
	// ReasonFollowedByFollowing 你关注的人关注了TA
	ReasonFollowedByFollowing ReasonType = iota
	// ReasonPopularInNetwork 在你的社交网络中很受欢迎
	ReasonPopularInNetwork
)

// RecommendationReason 值对象：推荐理由
//
// 这是一个复杂值对象的示例，展示了值对象不仅可以包装基本类型，
// 还可以封装复杂的业务概念和行为。
//
// 为什么推荐理由是值对象？
// 1. 推荐理由本身没有唯一标识
// 2. 两个相同类型、相同用户列表的理由是等价的
// 3. 推荐理由一旦创建就不应该改变
//
// 封装的业务规则：
// - 不同类型的推荐理由有不同的描述文案
// - 不同类型的推荐理由有不同的权重
// - 推荐理由的展示逻辑（Description 方法）
//
// 实际业务场景：
// 用户看到推荐时，会显示"3 位你关注的人也关注了TA"
// 这个文案的生成逻辑就封装在这个值对象中
//
// 重构说明（支持后端配置文案）：
// 新增 displayText 字段，用于存储后端返回的配置文案。
// 这样做的好处：
// 1. 后端可以动态配置文案（A/B测试、多语言、运营活动等）
// 2. 前端保留降级逻辑，后端接口异常时不影响用户体验
// 3. 渐进式迁移，不需要前后端同时上线
type RecommendationReason struct {
	reasonType   ReasonType
	relatedUsers []UserID // 哪些关注的人关注了这个推荐用户
	displayText  string   // 后端配置的展示文案（可选，为空时使用本地逻辑）
}

// NewFollowedByFollowingReason 工厂方法：创建"关注的人关注了TA"类型的推荐理由
func NewFollowedByFollowingReason(users []UserID) RecommendationReason {
	return RecommendationReason{
		reasonType:   ReasonFollowedByFollowing,
		relatedUsers: users,
		displayText:  "", // 使用本地逻辑生成文案
	}
}

// NewPopularInNetworkReason 工厂方法：创建"网络中受欢迎"类型的推荐理由
func NewPopularInNetworkReason(users []UserID) RecommendationReason {
	return RecommendationReason{
		reasonType:   ReasonPopularInNetwork,
		relatedUsers: users,
		displayText:  "", // 使用本地逻辑生成文案
	}
}

// NewRecommendationReasonWithText 工厂方法：创建带后端配置文案的推荐理由
//
// 这个工厂方法用于从后端接口数据创建推荐理由。
//
// 使用场景：
// 在 Application Service 或 DTO 转换层，将后端返回的数据映射到领域对象时使用。
//
// 示例：
//
//	// 后端返回的数据
//	apiResponse := {
//	    "reasonType": "followed_by_following",
//	    "displayText": "3 位你关注的人也关注了TA",
//	    "relatedUserIds": ["user1", "user2", "user3"]
//	}
//
//	// 在 Application Service 中转换
//	reason := NewRecommendationReasonWithText(
//	    ReasonFollowedByFollowing,
//	    userIDs,
//	    apiResponse.DisplayText,  // 使用后端配置的文案
//	)
//
// 设计思路：
// 1. 保持值对象的不可变性（通过工厂方法创建）
// 2. 明确区分"本地创建"和"从后端创建"两种场景
// 3. 为未来的扩展留出空间（如添加更多配置参数）
func NewRecommendationReasonWithText(reasonType ReasonType, users []UserID, displayText string) RecommendationReason {
	return RecommendationReason{
		reasonType:   reasonType,
		relatedUsers: users,
		displayText:  displayText, // 使用后端配置的文案
	}
}

// Description 生成用户可读的推荐理由描述
//
// 重构后的逻辑（渐进式迁移）：
// 1. 优先使用后端配置的文案（displayText）
// 2. 如果后端没有返回文案，降级到本地逻辑
//
// 为什么这样设计？
//
// 【优先使用后端配置】
// - 后端可以动态调整文案，无需发版
// - 支持 A/B 测试（不同用户看到不同文案）
// - 支持多语言（后端根据用户语言返回对应文案）
// - 支持运营活动（如节日特殊文案）
//
// 【保留本地降级逻辑】
// - 后端接口异常时，前端仍能正常展示
// - 兼容旧版本后端（还没有返回 displayText 的版本）
// - 灰度发布时，部分用户使用新逻辑，部分使用旧逻辑
// - 降低上线风险，可以随时回滚
//
// 实际场景：
//
//	// 场景1：后端返回了配置文案
//	reason1 := NewRecommendationReasonWithText(
//	    ReasonFollowedByFollowing,
//	    []UserID{user1, user2, user3},
//	    "你的 3 位好友也关注了TA", // 后端配置的新文案
//	)
//	desc1 := reason1.Description() // "你的 3 位好友也关注了TA"
//
//	// 场景2：后端还没有返回配置文案（兼容旧版本）
//	reason2 := NewFollowedByFollowingReason([]UserID{user1, user2, user3})
//	desc2 := reason2.Description() // "3 位你关注的人也关注了TA"（降级到本地逻辑）
//
// 迁移路径：
// 第1周：前端部署这个版本（支持 displayText，但保留降级逻辑）
// 第2周：后端灰度发布，开始返回 displayText
// 第3周：观察数据，确认稳定
// 第4周：后端全量发布
// 第N周：当确认后端稳定后，可以删除 switch 降级逻辑（可选）
//
// 注意事项：
// - displayText 为空字符串时，会使用本地逻辑
// - 后端应该保证返回的文案不为空，否则会降级
// - 如果需要强制使用后端文案（即使为空），可以增加一个标志位
func (r RecommendationReason) Description() string {
	// 优先使用后端配置的文案
	if r.displayText != "" {
		return r.displayText
	}

	// 降级到本地逻辑（兼容旧版本或后端异常）
	switch r.reasonType {
	case ReasonFollowedByFollowing:
		count := len(r.relatedUsers)
		if count == 1 {
			return "1 位你关注的人也关注了TA"
		}
		return fmt.Sprintf("%d 位你关注的人也关注了TA", count)
	case ReasonPopularInNetwork:
		return "在你的社交网络中很受欢迎"
	default:
		return "推荐给你"
	}
}

// RelatedUsers 访问器：获取相关用户列表
func (r RecommendationReason) RelatedUsers() []UserID {
	// 返回副本，保证不可变性
	result := make([]UserID, len(r.relatedUsers))
	copy(result, r.relatedUsers)
	return result
}

// Type 访问器：获取推荐理由类型
func (r RecommendationReason) Type() ReasonType {
	return r.reasonType
}

// Weight 业务规则：不同推荐理由的权重
//
// 这个方法展示了值对象如何参与业务计算。
//
// 业务规则：
// - 被更多人关注的用户权重更高（每个关注者 +10 分）
// - 不同类型的推荐理由有不同的基础权重
//
// 实际场景：
//
//	reason1 := NewFollowedByFollowingReason([]UserID{u1, u2, u3})
//	weight1 := reason1.Weight() // 3 × 10 = 30
//
//	reason2 := NewFollowedByFollowingReason([]UserID{u1})
//	weight2 := reason2.Weight() // 1 × 10 = 10
//
// 为什么权重计算在值对象中？
// 因为权重是推荐理由的固有属性，应该由值对象自己计算。
// 这样修改权重规则时，只需修改这一个地方。
//
// 扩展性：
// 未来可以添加更复杂的权重计算：
// - 考虑关注者的影响力
// - 考虑关注的时间衰减
// - 考虑用户的兴趣匹配度
func (r RecommendationReason) Weight() int {
	switch r.reasonType {
	case ReasonFollowedByFollowing:
		// 关注的人越多，权重越高
		return len(r.relatedUsers) * 10
	case ReasonPopularInNetwork:
		return 5
	default:
		return 1
	}
}
