package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"service/domain/entity"
	"service/domain/repository"
	"service/domain/valueobject"
)

// ContentRepositoryImpl 内容仓储实现
type ContentRepositoryImpl struct {
	db *gorm.DB
}

// NewContentRepository 构造函数
func NewContentRepository(db *gorm.DB) repository.ContentRepository {
	return &ContentRepositoryImpl{db: db}
}

// CountRecentPosts 实现接口：统计最近帖子数
func (r *ContentRepositoryImpl) CountRecentPosts(
	ctx context.Context,
	userID valueobject.UserID,
	days int,
) (int, error) {

	since := time.Now().AddDate(0, 0, -days)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&PostPO{}).
		Where("author_id = ? AND created_at >= ? AND status = ?",
			userID.Value(), since, "published").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetRecentPosts 实现接口：获取最近帖子
func (r *ContentRepositoryImpl) GetRecentPosts(
	ctx context.Context,
	userID valueobject.UserID,
	limit int,
) ([]*entity.Post, error) {

	var posts []PostPO
	err := r.db.WithContext(ctx).
		Where("author_id = ? AND status = ?", userID.Value(), "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	// 转换 PO -> 领域实体
	result := make([]*entity.Post, 0, len(posts))
	for _, po := range posts {
		postID, _ := valueobject.NewPostID(po.ID)
		authorID, _ := valueobject.NewUserID(po.AuthorID)

		post := entity.NewPost(postID, authorID, po.Content, po.CreatedAt)
		result = append(result, post)
	}

	return result, nil
}

// PostPO 帖子持久化对象
type PostPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	AuthorID  int64     `gorm:"index:idx_author;not null"`
	Content   string    `gorm:"type:text;not null"`
	Status    string    `gorm:"type:varchar(20);default:'published'"`
	CreatedAt time.Time `gorm:"index:idx_created_at;not null"`
	UpdatedAt time.Time
}

// TableName 指定表名
func (PostPO) TableName() string {
	return "posts"
}
