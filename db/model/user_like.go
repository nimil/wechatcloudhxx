package model

import "time"

// UserLikeModel 用户点赞关系模型
type UserLikeModel struct {
	Id        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    string    `gorm:"column:user_id;type:varchar(50);not null;index" json:"userId"`
	PostId    string    `gorm:"column:post_id;type:varchar(50);not null;index" json:"postId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName 指定表名
func (UserLikeModel) TableName() string {
	return "user_likes"
} 