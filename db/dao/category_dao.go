package dao

import (
	"wxcloudrun-golang/db/model"
)

// CategoryDao 分类数据访问接口
type CategoryDao interface {
	// 创建分类
	Create(category *model.CategoryModel) error
	
	// 根据ID获取分类
	GetById(id string) (*model.CategoryModel, error)
	
	// 根据代码获取分类
	GetByCode(code string) (*model.CategoryModel, error)
	
	// 获取所有分类
	GetAll() ([]*model.CategoryModel, error)
	
	// 获取可用于发布的分类
	GetForPublish() ([]*model.CategoryModel, error)
	
	// 更新分类
	Update(category *model.CategoryModel) error
	
	// 删除分类
	Delete(id string) error
	
	// 增加帖子数量
	IncrementPostCount(code string) error
	
	// 减少帖子数量
	DecrementPostCount(code string) error
} 