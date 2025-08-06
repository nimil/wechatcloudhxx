package model

import "time"

// CommentModel 评论模型
type CommentModel struct {
	Id        string    `gorm:"column:id;primaryKey;type:varchar(50)" json:"id"`
	Content   string    `gorm:"column:content;type:varchar(500);not null" json:"content"`
	AuthorId  string    `gorm:"column:author_id;type:varchar(50);not null;index" json:"authorId"`
	PostId    string    `gorm:"column:post_id;type:varchar(50);not null;index" json:"postId"`
	ParentId  *string   `gorm:"column:parent_id;type:varchar(50);index" json:"parentId"`
	Likes     int       `gorm:"column:likes;default:0" json:"likes"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (CommentModel) TableName() string {
	return "comments"
} 