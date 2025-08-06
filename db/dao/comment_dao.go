package dao

import (
	"wxcloudrun-golang/db/model"
)

// CommentDao 评论数据访问接口
type CommentDao interface {
	// 创建评论
	Create(comment *model.CommentModel) error
	
	// 根据ID获取评论
	GetById(id string) (*model.CommentModel, error)
	
	// 获取帖子评论列表
	GetByPostId(postId string, page, pageSize int) ([]*model.CommentModel, int64, error)
	
	// 更新评论
	Update(comment *model.CommentModel) error
	
	// 删除评论
	Delete(id string) error
	
	// 增加点赞数
	IncrementLikes(id string) error
	
	// 减少点赞数
	DecrementLikes(id string) error
} 