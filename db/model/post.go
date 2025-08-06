package model

import "time"

// PostModel 帖子模型
type PostModel struct {
	Id           string    `gorm:"column:id;primaryKey;type:varchar(50)" json:"id"`
	Title        string    `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Content      string    `gorm:"column:content;type:text;not null" json:"content"`
	Excerpt      string    `gorm:"column:excerpt;type:varchar(500)" json:"excerpt"`
	AuthorId     string    `gorm:"column:author_id;type:varchar(50);not null;index" json:"authorId"`
	Category     string    `gorm:"column:category;type:varchar(20);not null;index" json:"category"`
	CategoryName string    `gorm:"column:category_name;type:varchar(50);not null" json:"categoryName"`
	Tags         string    `gorm:"column:tags;type:text" json:"tags"` // JSON格式存储
	Images       string    `gorm:"column:images;type:text" json:"images"` // JSON格式存储
	IsPublic     bool      `gorm:"column:is_public;default:true" json:"isPublic"`
	Likes        int       `gorm:"column:likes;default:0" json:"likes"`
	Comments     int       `gorm:"column:comments;default:0" json:"comments"`
	Views        int       `gorm:"column:views;default:0" json:"views"`
	Shares       int       `gorm:"column:shares;default:0" json:"shares"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (PostModel) TableName() string {
	return "posts"
} 