package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CreateUser 创建用户
func (imp *UserInterfaceImp) CreateUser(user *model.UserModel) error {
	cli := db.Get()
	return cli.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func (imp *UserInterfaceImp) GetUserByUsername(username string) (*model.UserModel, error) {
	var err error
	var user = new(model.UserModel)

	cli := db.Get()
	err = cli.Where("username = ?", username).First(user).Error

	return user, err
}

// GetUserById 根据ID查询用户
func (imp *UserInterfaceImp) GetUserById(id int32) (*model.UserModel, error) {
	var err error
	var user = new(model.UserModel)

	cli := db.Get()
	err = cli.Where("id = ?", id).First(user).Error

	return user, err
}

// GetUsersByPage 分页查询用户列表
func (imp *UserInterfaceImp) GetUsersByPage(page, pageSize int) ([]*model.UserModel, int64, error) {
	var users []*model.UserModel
	var total int64

	cli := db.Get()

	// 获取总数
	err := cli.Model(&model.UserModel{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = cli.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
