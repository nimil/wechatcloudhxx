package dao

import (
	"wxcloudrun-golang/db/model"
)

// PostDao 帖子数据访问接口
type PostDao interface {
	// 创建帖子
	Create(post *model.PostModel) error
	
	// 根据ID获取帖子
	GetById(id int64) (*model.PostModel, error)
	
	// 获取帖子列表
	GetList(page, pageSize int, category, sort string) ([]*model.PostModel, int64, error)
	
	// 更新帖子
	Update(post *model.PostModel) error
	
	// 删除帖子
	Delete(id int64) error
	
	// 增加浏览量
	IncrementViews(id int64) error
	
	// 增加点赞数
	IncrementLikes(id int64) error
	
	// 减少点赞数
	DecrementLikes(id int64) error
	
	// 增加评论数
	IncrementComments(id int64) error
	
	// 减少评论数
	DecrementComments(id int64) error
} 