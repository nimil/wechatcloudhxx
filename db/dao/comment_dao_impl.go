package dao

import (
	"gorm.io/gorm"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CommentDaoImpl 评论DAO实现
type CommentDaoImpl struct {
	db *gorm.DB
}

// NewCommentDao 创建评论DAO实例
func NewCommentDao() CommentDao {
	return &CommentDaoImpl{db: db.GetDB()}
}

// Create 创建评论
func (dao *CommentDaoImpl) Create(comment *model.CommentModel) error {
	return dao.db.Create(comment).Error
}

// GetById 根据ID获取评论
func (dao *CommentDaoImpl) GetById(id string) (*model.CommentModel, error) {
	var comment model.CommentModel
	err := dao.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByPostId 获取帖子评论列表
func (dao *CommentDaoImpl) GetByPostId(postId string, page, pageSize int) ([]*model.CommentModel, int64, error) {
	var comments []*model.CommentModel
	var total int64
	
	query := dao.db.Model(&model.CommentModel{}).Where("post_id = ? AND parent_id IS NULL", postId)
	
	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 分页
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}
	
	return comments, total, nil
}

// Update 更新评论
func (dao *CommentDaoImpl) Update(comment *model.CommentModel) error {
	return dao.db.Save(comment).Error
}

// Delete 删除评论
func (dao *CommentDaoImpl) Delete(id string) error {
	return dao.db.Where("id = ?", id).Delete(&model.CommentModel{}).Error
}

// IncrementLikes 增加点赞数
func (dao *CommentDaoImpl) IncrementLikes(id string) error {
	return dao.db.Model(&model.CommentModel{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}

// DecrementLikes 减少点赞数
func (dao *CommentDaoImpl) DecrementLikes(id string) error {
	return dao.db.Model(&model.CommentModel{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
} 