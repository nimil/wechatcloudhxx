package dao

import "wxcloudrun-golang/db/model"

// ImageCheckDao 图片检测数据访问接口
type ImageCheckDao interface {
	// Create 创建图片检测记录
	Create(imageCheck *model.ImageCheckModel) error
	
	// GetByTraceId 根据trace_id获取检测记录
	GetByTraceId(traceId string) (*model.ImageCheckModel, error)
	
	// UpdateStatus 更新检测状态
	UpdateStatus(traceId string, status int, suggest string, label int, prob float64, strategy string, errcode int, errmsg string) error
	
	// GetByPostId 获取帖子的所有图片检测记录
	GetByPostId(postId int64) ([]*model.ImageCheckModel, error)
	
	// GetPendingChecks 获取待检测的记录
	GetPendingChecks() ([]*model.ImageCheckModel, error)
	
	// DeleteByPostId 删除帖子的所有检测记录
	DeleteByPostId(postId int64) error
}
