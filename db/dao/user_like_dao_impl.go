package dao

import (
	"gorm.io/gorm"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// UserLikeDaoImpl 用户点赞DAO实现
type UserLikeDaoImpl struct {
	db *gorm.DB
}

// NewUserLikeDao 创建用户点赞DAO实例
func NewUserLikeDao() UserLikeDao {
	return &UserLikeDaoImpl{db: db.GetDB()}
}

// Create 创建点赞记录
func (dao *UserLikeDaoImpl) Create(userLike *model.UserLikeModel) error {
	return dao.db.Create(userLike).Error
}

// Delete 删除点赞记录
func (dao *UserLikeDaoImpl) Delete(userId, postId int64) error {
	return dao.db.Where("user_id = ? AND post_id = ?", userId, postId).Delete(&model.UserLikeModel{}).Error
}

// IsLiked 检查用户是否点赞
func (dao *UserLikeDaoImpl) IsLiked(userId, postId int64) (bool, error) {
	var count int64
	err := dao.db.Model(&model.UserLikeModel{}).Where("user_id = ? AND post_id = ?", userId, postId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserLikedPostIds 获取用户点赞的帖子ID列表
func (dao *UserLikeDaoImpl) GetUserLikedPostIds(userId int64) ([]int64, error) {
	var postIds []int64
	err := dao.db.Model(&model.UserLikeModel{}).Where("user_id = ?", userId).Pluck("post_id", &postIds).Error
	if err != nil {
		return nil, err
	}
	return postIds, nil
} 