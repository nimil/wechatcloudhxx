package model

import "time"

// UserModel 用户模型
type UserModel struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"column:password;not null" json:"password"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}
