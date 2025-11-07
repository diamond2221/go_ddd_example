package valueobject

import "errors"

var (
	ErrInvalidPostID = errors.New("invalid post id: must be positive")
)

// PostID 值对象：帖子ID
type PostID struct {
	value int64
}

// NewPostID 工厂方法
func NewPostID(value int64) (PostID, error) {
	if value <= 0 {
		return PostID{}, ErrInvalidPostID
	}
	return PostID{value: value}, nil
}

// Value 访问器
func (p PostID) Value() int64 {
	return p.value
}

// Equals 值对象相等性比较
func (p PostID) Equals(other PostID) bool {
	return p.value == other.value
}

type D struct {
	value int64
}
