package dao

import (
	"wxcloudrun-golang/db/model"
)

// UserDao 用户数据访问接口
type UserDao interface {
	CreateUser(user *model.UserModel) error
	GetUserByUsername(username string) (*model.UserModel, error)
	GetById(id string) (*model.UserModel, error)
	GetUserByOpenId(openId string) (*model.UserModel, error)
	GetUserByUnionId(unionId string) (*model.UserModel, error)
	GetUsersByPage(page, pageSize int) ([]*model.UserModel, int64, error)
	UpdateUser(user *model.UserModel) error
	DeleteUser(id string) error
}

// UserDaoImpl 用户数据访问实现
type UserDaoImpl struct{}

// NewUserDao 创建用户DAO实例
func NewUserDao() UserDao {
	return &UserDaoImpl{}
}
