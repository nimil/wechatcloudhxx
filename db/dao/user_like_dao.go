package dao

import (
	"wxcloudrun-golang/db/model"
)

// UserLikeDao 用户点赞数据访问接口
type UserLikeDao interface {
	// 创建点赞记录
	Create(userLike *model.UserLikeModel) error
	
	// 删除点赞记录
	Delete(userId, postId int64) error
	
	// 检查用户是否点赞
	IsLiked(userId, postId int64) (bool, error)
	
	// 获取用户点赞的帖子ID列表
	GetUserLikedPostIds(userId int64) ([]int64, error)
} 