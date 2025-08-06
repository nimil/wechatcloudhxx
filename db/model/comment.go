package model

import "time"

// CommentModel 评论模型
type CommentModel struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"column:content;type:varchar(500);not null" json:"content"`
	AuthorId  int64     `gorm:"column:author_id;not null;index" json:"authorId"`
	PostId    int64     `gorm:"column:post_id;not null;index" json:"postId"`
	ParentId  *int64    `gorm:"column:parent_id;index" json:"parentId"`
	Likes     int       `gorm:"column:likes;default:0" json:"likes"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (CommentModel) TableName() string {
	return "comments"
} 