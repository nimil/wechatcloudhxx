package model

import "time"

// CategoryModel 分类模型
type CategoryModel struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Code        string    `gorm:"column:code;type:varchar(20);not null;uniqueIndex" json:"code"`
	Icon        string    `gorm:"column:icon;type:varchar(50)" json:"icon"`
	Description string    `gorm:"column:description;type:varchar(200)" json:"description"`
	PostCount   int       `gorm:"column:post_count;default:0" json:"postCount"`
	Sort        int       `gorm:"column:sort;default:0" json:"sort"`
	IsActive    bool      `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (CategoryModel) TableName() string {
	return "categories"
} 