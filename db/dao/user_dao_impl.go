package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CreateUser 创建用户
func (dao *UserDaoImpl) CreateUser(user *model.UserModel) error {
	return db.GetDB().Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func (dao *UserDaoImpl) GetUserByUsername(username string) (*model.UserModel, error) {
	var user model.UserModel
	err := db.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetById 根据ID查询用户
func (dao *UserDaoImpl) GetById(id string) (*model.UserModel, error) {
	var user model.UserModel
	err := db.GetDB().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByOpenId 根据OpenId查询用户
func (dao *UserDaoImpl) GetUserByOpenId(openId string) (*model.UserModel, error) {
	var user model.UserModel
	err := db.GetDB().Where("openid = ?", openId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUnionId 根据UnionId查询用户
func (dao *UserDaoImpl) GetUserByUnionId(unionId string) (*model.UserModel, error) {
	var user model.UserModel
	err := db.GetDB().Where("unionid = ?", unionId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByPage 分页查询用户列表
func (dao *UserDaoImpl) GetUsersByPage(page, pageSize int) ([]*model.UserModel, int64, error) {
	var users []*model.UserModel
	var total int64

	// 获取总数
	err := db.GetDB().Model(&model.UserModel{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = db.GetDB().Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser 更新用户信息
func (dao *UserDaoImpl) UpdateUser(user *model.UserModel) error {
	return db.GetDB().Save(user).Error
}

// DeleteUser 删除用户
func (dao *UserDaoImpl) DeleteUser(id string) error {
	return db.GetDB().Where("id = ?", id).Delete(&model.UserModel{}).Error
} 