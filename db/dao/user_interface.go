package dao

import (
	"wxcloudrun-golang/db/model"
)

// UserInterface 用户数据模型接口
type UserInterface interface {
	CreateUser(user *model.UserModel) error
	GetUserByUsername(username string) (*model.UserModel, error)
	GetUserById(id int32) (*model.UserModel, error)
	GetUsersByPage(page, pageSize int) ([]*model.UserModel, int64, error)
}

// UserInterfaceImp 用户数据模型实现
type UserInterfaceImp struct{}

// UserImp 用户实现实例
var UserImp UserInterface = &UserInterfaceImp{}
