package valueobject

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	ErrNicknameTooShort      = errors.New("昵称长度不能少于3个字符")
	ErrNicknameTooLong       = errors.New("昵称长度不能超过16个字符")
	ErrNicknameInvalidFormat = errors.New("昵称只能包含中文、英文字母和数字")
)

// Nickname 值对象：用户昵称
//
// 为什么昵称需要定义成值对象？
// 1. 有明确的业务规则（长度限制、字符限制）
// 2. 验证逻辑复杂，需要集中管理
// 3. 昵称在多处使用，避免到处重复验证
// 4. 保证数据有效性，创建后就是合法的昵称
//
// 业务规则：
// - 最短 3 个字符，最长 16 个字符
// - 只能包含中文、英文字母（大小写）、数字
// - 不可变：创建后不能修改
//
// 使用示例：
//
//	// 创建昵称
//	nickname, err := NewNickname("张三123")
//	if err != nil {
//	    // 处理验证失败
//	}
//
//	// 使用昵称
//	user := domain.NewUser(userID, nickname, email)
type Nickname struct {
	value string
}

// 正则表达式：只允许中文、英文字母、数字
// \x{4e00}-\x{9fa5} 匹配中文字符（使用 Go 的 Unicode 范围语法）
// a-zA-Z 匹配英文字母
// 0-9 匹配数字
var nicknamePattern = regexp.MustCompile(`^[\p{Han}a-zA-Z0-9]+$`)

// NewNickname 工厂方法：创建昵称值对象
//
// 在创建时验证所有业务规则：
// 1. 长度检查（3-16个字符）
// 2. 字符格式检查（中文/英文/数字）
//
// 注意：这里使用 utf8.RuneCountInString 而不是 len()
// 因为中文字符占多个字节，len() 会得到错误的长度
//
// 示例：
//
//	nickname1, _ := NewNickname("张三")           // ❌ 太短
//	nickname2, _ := NewNickname("张三123")        // ✅ 合法
//	nickname3, _ := NewNickname("Alice")          // ✅ 合法
//	nickname4, _ := NewNickname("用户123")        // ✅ 合法
//	nickname5, _ := NewNickname("张三@123")       // ❌ 包含特殊字符
//	nickname6, _ := NewNickname("这是一个超级超级长的昵称") // ❌ 太长
func NewNickname(value string) (Nickname, error) {
	// 规则1：长度检查（使用字符数而不是字节数）
	length := utf8.RuneCountInString(value)
	if length < 3 {
		return Nickname{}, ErrNicknameTooShort
	}
	if length > 16 {
		return Nickname{}, ErrNicknameTooLong
	}

	// 规则2：字符格式检查（只允许中文、英文、数字）
	if !nicknamePattern.MatchString(value) {
		return Nickname{}, ErrNicknameInvalidFormat
	}

	return Nickname{value: value}, nil
}

// Value 访问器：获取昵称字符串
func (n Nickname) Value() string {
	return n.value
}

// Equals 值对象相等性比较
func (n Nickname) Equals(other Nickname) bool {
	return n.value == other.value
}

// String 实现 Stringer 接口
func (n Nickname) String() string {
	return n.value
}

// Length 获取昵称长度（字符数）
//
// 这个方法展示了值对象可以提供便捷的业务方法
// 外部调用者不需要关心如何正确计算中文字符长度
func (n Nickname) Length() int {
	return utf8.RuneCountInString(n.value)
}
