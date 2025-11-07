package valueobject

import (
	"github.com/google/uuid"
)

// RecommendationID 值对象：推荐ID
// 使用 UUID 作为唯一标识
type RecommendationID struct {
	value string
}

// NewRecommendationID 工厂方法：生成新的推荐ID
func NewRecommendationID() RecommendationID {
	return RecommendationID{
		value: uuid.New().String(),
	}
}

// FromString 工厂方法：从字符串创建推荐ID
func RecommendationIDFromString(value string) (RecommendationID, error) {
	// 验证是否是有效的 UUID
	if _, err := uuid.Parse(value); err != nil {
		return RecommendationID{}, err
	}
	return RecommendationID{value: value}, nil
}

// Value 访问器
func (r RecommendationID) Value() string {
	return r.value
}

// Equals 值对象相等性比较
func (r RecommendationID) Equals(other RecommendationID) bool {
	return r.value == other.value
}

// String 实现 Stringer 接口
func (r RecommendationID) String() string {
	return r.value
}
