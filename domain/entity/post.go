package entity

import (
	"time"

	"service/domain/valueobject"
)

// Post 实体：帖子
//
// 什么是实体？
// 实体是有唯一标识的领域对象，即使属性相同，只要 ID 不同就是不同的对象。
//
// 实体 vs 值对象：
// ┌──────────┬────────────────┬────────────────┐
// │          │ 实体（Post）    │ 值对象（UserID）│
// ├──────────┼────────────────┼────────────────┤
// │ 标识     │ 有（PostID）    │ 无              │
// │ 可变性   │ 可变            │ 不可变          │
// │ 相等性   │ 通过ID比较      │ 通过值比较      │
// │ 生命周期 │ 有              │ 无              │
// └──────────┴────────────────┴────────────────┘
//
// 实际示例：
//
//	post1 := NewPost(PostID(1), userA, "Hello", time.Now())
//	post2 := NewPost(PostID(1), userA, "World", time.Now())
//	post1.Equals(post2) // true，因为 ID 相同
//
//	userID1 := NewUserID(123)
//	userID2 := NewUserID(123)
//	userID1.Equals(userID2) // true，因为值相同
//
// 实体 vs 聚合根：
// - Post 在推荐上下文中是简单实体，不是聚合根
// - 它不需要维护复杂的一致性边界
// - 如果在内容管理上下文中，Post 可能是聚合根（管理评论、点赞等）
//
// 上下文边界（Bounded Context）：
// 同一个概念在不同上下文中可能有不同的角色：
// - 推荐上下文：Post 是简单实体，只关心内容和作者
// - 内容上下文：Post 是聚合根，管理评论、点赞、审核状态等
type Post struct {
	id        valueobject.PostID
	authorID  valueobject.UserID
	content   string
	createdAt time.Time
}

// NewPost 工厂方法
func NewPost(
	id valueobject.PostID,
	authorID valueobject.UserID,
	content string,
	createdAt time.Time,
) *Post {
	return &Post{
		id:        id,
		authorID:  authorID,
		content:   content,
		createdAt: createdAt,
	}
}

// --- 访问器方法 ---

func (p *Post) ID() valueobject.PostID {
	return p.id
}

func (p *Post) AuthorID() valueobject.UserID {
	return p.authorID
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) CreatedAt() time.Time {
	return p.createdAt
}
