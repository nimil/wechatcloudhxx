package model

import "time"

// UserLikeModel 用户点赞关系模型
type UserLikeModel struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;not null;index" json:"userId"`
	PostId    int64     `gorm:"column:post_id;not null;index" json:"postId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName 指定表名
func (UserLikeModel) TableName() string {
	return "user_likes"
} 