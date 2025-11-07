package valueobject

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUserID = errors.New("invalid user id: must be positive")
)

// UserID 值对象：用户ID
//
// 什么是值对象？
// 值对象是 DDD 中的核心概念，用于封装业务概念和验证规则。
//
// 值对象的三大特征：
// 1. 不可变性：创建后不能修改，所有字段私有
// 2. 无唯一标识：通过值比较相等性，而不是通过 ID
// 3. 自我验证：在创建时就保证数据有效性
//
// 为什么不直接用 int64？
// 传统方式：func GetUser(userID int64) - 容易传错参数
// 值对象方式：func GetUser(userID UserID) - 类型安全，语义清晰
//
// 实际好处：
// - 类型安全：不会把 postID 误传给需要 userID 的函数
// - 业务规则集中：所有关于 UserID 的验证都在这里
// - 代码可读性：看到 UserID 就知道是用户标识
//
// 使用示例：
//
//	userID, err := valueobject.NewUserID(123)
//	if err != nil {
//	    // 处理无效ID
//	}
//	// userID 保证是有效的，可以安全使用
type UserID struct {
	value int64 // 私有字段，外部无法直接访问和修改
}

// NewUserID 工厂方法（构造函数）
//
// 为什么用工厂方法？
// 1. 验证：在创建时就保证数据有效性
// 2. 封装：隐藏内部实现细节
// 3. 灵活：未来可以改变创建逻辑而不影响调用者
//
// 在创建时验证业务规则：
// - 用户ID必须是正数（业务规则）
// - 如果验证失败，返回错误而不是创建无效对象
//
// 使用示例：
//
//	userID, err := NewUserID(123)
//	if err != nil {
//	    // 处理无效ID
//	}
//	// userID 保证是有效的
//
// 对比直接创建：
//
//	userID := UserID{value: -1} // 编译通过，但数据无效
//	工厂方法保证了数据的有效性
func NewUserID(value int64) (UserID, error) {
	// 业务规则：用户ID必须是正数
	if value <= 0 {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: value}, nil
}

// Value 访问器方法
// 只读访问，保证不可变性
func (u UserID) Value() int64 {
	return u.value
}

// Equals 值对象通过值比较相等性
// 两个 UserID 只要值相同就相等
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// String 实现 Stringer 接口，方便日志输出
func (u UserID) String() string {
	return fmt.Sprintf("UserID(%d)", u.value)
}
