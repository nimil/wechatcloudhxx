package service

import (
	"fmt"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// LikeService 点赞服务
type LikeService struct {
	userLikeDao dao.UserLikeDao
	postDao     dao.PostDao
}

// NewLikeService 创建点赞服务实例
func NewLikeService() *LikeService {
	return &LikeService{
		userLikeDao: dao.NewUserLikeDao(),
		postDao:     dao.NewPostDao(),
	}
}

// LikeRequest 点赞请求
type LikeRequest struct {
	Action string `json:"action"` // like 或 unlike
}

// LikeResponse 点赞响应
type LikeResponse struct {
	IsLiked    bool `json:"isLiked"`
	LikesCount int  `json:"likesCount"`
}

// ToggleLike 切换点赞状态
func (s *LikeService) ToggleLike(postId string, userId string, req *LikeRequest) (*LikeResponse, error) {
	// 验证帖子是否存在
	_, err := s.postDao.GetById(postId)
	if err != nil {
		return nil, fmt.Errorf("帖子不存在: %v", err)
	}

	// 检查当前点赞状态
	isLiked, err := s.userLikeDao.IsLiked(userId, postId)
	if err != nil {
		return nil, fmt.Errorf("检查点赞状态失败: %v", err)
	}

	switch req.Action {
	case "like":
		if !isLiked {
			// 创建点赞记录
			userLike := &model.UserLikeModel{
				UserId: userId,
				PostId: postId,
			}
			err = s.userLikeDao.Create(userLike)
			if err != nil {
				return nil, fmt.Errorf("创建点赞记录失败: %v", err)
			}

			// 增加帖子点赞数
			err = s.postDao.IncrementLikes(postId)
			if err != nil {
				return nil, fmt.Errorf("更新帖子点赞数失败: %v", err)
			}

			isLiked = true
		}

	case "unlike":
		if isLiked {
			// 删除点赞记录
			err = s.userLikeDao.Delete(userId, postId)
			if err != nil {
				return nil, fmt.Errorf("删除点赞记录失败: %v", err)
			}

			// 减少帖子点赞数
			err = s.postDao.DecrementLikes(postId)
			if err != nil {
				return nil, fmt.Errorf("更新帖子点赞数失败: %v", err)
			}

			isLiked = false
		}

	default:
		return nil, fmt.Errorf("无效的操作: %s", req.Action)
	}

	// 获取最新的点赞数
	updatedPost, err := s.postDao.GetById(postId)
	if err != nil {
		return nil, fmt.Errorf("获取帖子信息失败: %v", err)
	}

	return &LikeResponse{
		IsLiked:    isLiked,
		LikesCount: updatedPost.Likes,
	}, nil
} 