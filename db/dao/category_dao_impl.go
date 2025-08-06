package dao

import (
	"gorm.io/gorm"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CategoryDaoImpl 分类DAO实现
type CategoryDaoImpl struct {
	db *gorm.DB
}

// NewCategoryDao 创建分类DAO实例
func NewCategoryDao() CategoryDao {
	return &CategoryDaoImpl{db: db.GetDB()}
}

// Create 创建分类
func (dao *CategoryDaoImpl) Create(category *model.CategoryModel) error {
	return dao.db.Create(category).Error
}

// GetById 根据ID获取分类
func (dao *CategoryDaoImpl) GetById(id string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	err := dao.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByCode 根据代码获取分类
func (dao *CategoryDaoImpl) GetByCode(code string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	err := dao.db.Where("code = ?", code).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类
func (dao *CategoryDaoImpl) GetAll() ([]*model.CategoryModel, error) {
	var categories []*model.CategoryModel
	err := dao.db.Where("is_active = ?", true).Order("sort ASC, post_count DESC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetForPublish 获取可用于发布的分类
func (dao *CategoryDaoImpl) GetForPublish() ([]*model.CategoryModel, error) {
	var categories []*model.CategoryModel
	err := dao.db.Where("is_active = ? AND code != ?", true, "all").Order("sort ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 更新分类
func (dao *CategoryDaoImpl) Update(category *model.CategoryModel) error {
	return dao.db.Save(category).Error
}

// Delete 删除分类
func (dao *CategoryDaoImpl) Delete(id string) error {
	return dao.db.Where("id = ?", id).Delete(&model.CategoryModel{}).Error
}

// IncrementPostCount 增加帖子数量
func (dao *CategoryDaoImpl) IncrementPostCount(code string) error {
	return dao.db.Model(&model.CategoryModel{}).Where("code = ?", code).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// DecrementPostCount 减少帖子数量
func (dao *CategoryDaoImpl) DecrementPostCount(code string) error {
	return dao.db.Model(&model.CategoryModel{}).Where("code = ?", code).UpdateColumn("post_count", gorm.Expr("post_count - ?", 1)).Error
} 