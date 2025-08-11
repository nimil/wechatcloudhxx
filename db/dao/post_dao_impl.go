package dao

import (
	"gorm.io/gorm"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// PostDaoImpl 帖子DAO实现
type PostDaoImpl struct {
	db *gorm.DB
}

// NewPostDao 创建帖子DAO实例
func NewPostDao() PostDao {
	return &PostDaoImpl{db: db.GetDB()}
}

// Create 创建帖子
func (dao *PostDaoImpl) Create(post *model.PostModel) error {
	return dao.db.Create(post).Error
}

// GetById 根据ID获取帖子
func (dao *PostDaoImpl) GetById(id int64) (*model.PostModel, error) {
	var post model.PostModel
	err := dao.db.Where("id = ? AND is_deleted = ?", id, false).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetList 获取帖子列表
func (dao *PostDaoImpl) GetList(page, pageSize int, category, sort string) ([]*model.PostModel, int64, error) {
	var posts []*model.PostModel
	var total int64
	
	query := dao.db.Model(&model.PostModel{}).Where("is_public = ? AND is_deleted = ?", true, false)
	
	// 分类筛选
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}
	
	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 排序
	switch sort {
	case "hot":
		query = query.Order("likes DESC, views DESC, created_at DESC")
	case "recommend":
		query = query.Order("likes DESC, comments DESC, created_at DESC")
	default: // latest
		query = query.Order("created_at DESC")
	}
	
	// 分页
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}
	
	return posts, total, nil
}

// Update 更新帖子
func (dao *PostDaoImpl) Update(post *model.PostModel) error {
	return dao.db.Save(post).Error
}

// Delete 删除帖子（物理删除）
func (dao *PostDaoImpl) Delete(id int64) error {
	return dao.db.Where("id = ?", id).Delete(&model.PostModel{}).Error
}

// SoftDelete 逻辑删除帖子
func (dao *PostDaoImpl) SoftDelete(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// Restore 恢复帖子
func (dao *PostDaoImpl) Restore(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).Update("is_deleted", false).Error
}

// IncrementViews 增加浏览量
func (dao *PostDaoImpl) IncrementViews(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// IncrementLikes 增加点赞数
func (dao *PostDaoImpl) IncrementLikes(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}

// DecrementLikes 减少点赞数
func (dao *PostDaoImpl) DecrementLikes(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
}

// IncrementComments 增加评论数
func (dao *PostDaoImpl) IncrementComments(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).UpdateColumn("comments", gorm.Expr("comments + ?", 1)).Error
}

// DecrementComments 减少评论数
func (dao *PostDaoImpl) DecrementComments(id int64) error {
	return dao.db.Model(&model.PostModel{}).Where("id = ?", id).UpdateColumn("comments", gorm.Expr("comments - ?", 1)).Error
} 