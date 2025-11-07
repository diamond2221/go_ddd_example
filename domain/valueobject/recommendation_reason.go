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
type RecommendationReason struct {
	reasonType   ReasonType
	relatedUsers []UserID // 哪些关注的人关注了这个推荐用户
}

// NewFollowedByFollowingReason 工厂方法：创建"关注的人关注了TA"类型的推荐理由
func NewFollowedByFollowingReason(users []UserID) RecommendationReason {
	return RecommendationReason{
		reasonType:   ReasonFollowedByFollowing,
		relatedUsers: users,
	}
}

// NewPopularInNetworkReason 工厂方法：创建"网络中受欢迎"类型的推荐理由
func NewPopularInNetworkReason(users []UserID) RecommendationReason {
	return RecommendationReason{
		reasonType:   ReasonPopularInNetwork,
		relatedUsers: users,
	}
}

// Description 生成用户可读的推荐理由描述
//
// 这个方法展示了值对象如何封装展示逻辑。
//
// 为什么展示逻辑在领域层？
// 因为"如何描述推荐理由"是业务规则，不是技术细节。
// 产品经理定义了文案规则：
// - 1个人关注：显示"1 位你关注的人也关注了TA"
// - 多个人关注：显示"N 位你关注的人也关注了TA"
//
// 实际场景：
//
//	reason := NewFollowedByFollowingReason([]UserID{user1, user2, user3})
//	desc := reason.Description() // "3 位你关注的人也关注了TA"
//
// 国际化考虑：
// 如果需要支持多语言，可以：
// 1. 返回 key，由外层翻译
// 2. 传入 locale 参数
// 3. 使用 i18n 库
//
// 为什么不在前端生成文案？
// 因为文案规则是业务规则，应该由后端控制。
// 前端只负责展示，不应该包含业务逻辑。
func (r RecommendationReason) Description() string {
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
