package model

import "time"

// UserModel 用户模型
type UserModel struct {
	Id          string    `gorm:"column:id;primaryKey;type:varchar(50)" json:"id"`
	Username    string    `gorm:"column:username;uniqueIndex;not null" json:"username"`
	Nickname    string    `gorm:"column:nickname;type:varchar(50)" json:"nickname"`
	Avatar      string    `gorm:"column:avatar;type:varchar(500)" json:"avatar"`
	Bio         string    `gorm:"column:bio;type:varchar(200)" json:"bio"`
	Level       int       `gorm:"column:level;default:1" json:"level"`
	IsVerified  bool      `gorm:"column:is_verified;default:false" json:"isVerified"`
	Password    string    `gorm:"column:password;not null" json:"password"`
	OpenId      string    `gorm:"column:openid;index" json:"openid"`
	UnionId     string    `gorm:"column:unionid;index" json:"unionid"`
	AppId       string    `gorm:"column:appid" json:"appid"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}
