package service

import (
	"wxcloudrun-golang/db/dao"
)

// CategoryService 分类服务
type CategoryService struct {
	categoryDao dao.CategoryDao
}

// NewCategoryService 创建分类服务实例
func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryDao: dao.NewCategoryDao(),
	}
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	PostCount   int    `json:"postCount"`
}

// TopicInfo 话题信息
type TopicInfo struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Code        string `json:"code"`
	PostCount   int    `json:"postCount"`
	FollowCount int    `json:"followCount"`
	IsFollowed  bool   `json:"isFollowed"`
}

// GetCategories 获取所有分类
func (s *CategoryService) GetCategories() ([]*CategoryInfo, error) {
	categories, err := s.categoryDao.GetAll()
	if err != nil {
		return nil, err
	}

	categoryInfos := make([]*CategoryInfo, 0, len(categories))
	for _, category := range categories {
		categoryInfo := &CategoryInfo{
			Id:          category.Id,
			Name:        category.Name,
			Code:        category.Code,
			Icon:        category.Icon,
			Description: category.Description,
			PostCount:   category.PostCount,
		}
		categoryInfos = append(categoryInfos, categoryInfo)
	}

	return categoryInfos, nil
}

// GetPublishCategories 获取可用于发布的分类
func (s *CategoryService) GetPublishCategories() ([]*CategoryInfo, error) {
	categories, err := s.categoryDao.GetForPublish()
	if err != nil {
		return nil, err
	}

	categoryInfos := make([]*CategoryInfo, 0, len(categories))
	for _, category := range categories {
		categoryInfo := &CategoryInfo{
			Id:          category.Id,
			Name:        category.Name,
			Code:        category.Code,
			Icon:        category.Icon,
			Description: category.Description,
			PostCount:   category.PostCount,
		}
		categoryInfos = append(categoryInfos, categoryInfo)
	}

	return categoryInfos, nil
}

// GetHotTopics 获取热门话题
func (s *CategoryService) GetHotTopics(userId int64) ([]*TopicInfo, error) {
	categories, err := s.categoryDao.GetAll()
	if err != nil {
		return nil, err
	}

	// 这里简化处理，实际应该根据帖子数量和关注数来排序
	topicInfos := make([]*TopicInfo, 0, len(categories))
	for _, category := range categories {
		// 跳过"全部"分类
		if category.Code == "all" {
			continue
		}

		topicInfo := &TopicInfo{
			Id:          category.Id,
			Name:        category.Name,
			Icon:        category.Icon,
			Code:        category.Code,
			PostCount:   category.PostCount,
			FollowCount: category.PostCount / 2, // 简化处理，实际应该从数据库获取
			IsFollowed:  false,                   // TODO: 实现关注功能
		}
		topicInfos = append(topicInfos, topicInfo)
	}

	return topicInfos, nil
} 