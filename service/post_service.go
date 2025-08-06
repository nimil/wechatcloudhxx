package service

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// PostService 帖子服务
type PostService struct {
	postDao     dao.PostDao
	userDao     dao.UserDao
	categoryDao dao.CategoryDao
	userLikeDao dao.UserLikeDao
}

// NewPostService 创建帖子服务实例
func NewPostService() *PostService {
	return &PostService{
		postDao:     dao.NewPostDao(),
		userDao:     dao.NewUserDao(),
		categoryDao: dao.NewCategoryDao(),
		userLikeDao: dao.NewUserLikeDao(),
	}
}

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Images   []string `json:"images"`
	IsPublic bool     `json:"isPublic"`
}

// CreatePostResponse 创建帖子响应
type CreatePostResponse struct {
	PostId    int64     `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
	URL       string    `json:"url"`
}

// PostListResponse 帖子列表响应
type PostListResponse struct {
	List       []*PostDetail `json:"list"`
	Pagination Pagination    `json:"pagination"`
}

// PostDetail 帖子详情
type PostDetail struct {
	Id           int64     `json:"id"`
	Title        string    `json:"title"`
	Excerpt      string    `json:"excerpt"`
	Content      string    `json:"content"`
	Author       UserInfo  `json:"author"`
	Category     string    `json:"category"`
	CategoryName string    `json:"categoryName"`
	Tags         []string  `json:"tags"`
	Images       []string  `json:"images"`
	Stats        PostStats `json:"stats"`
	IsLiked      bool      `json:"isLiked"`
	IsCollected  bool      `json:"isCollected"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UserInfo 用户信息
type UserInfo struct {
	Id         int64  `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Bio        string `json:"bio"`
	Level      int    `json:"level"`
	IsVerified bool   `json:"isVerified"`
}

// PostStats 帖子统计
type PostStats struct {
	Likes    int `json:"likes"`
	Comments int `json:"comments"`
	Views    int `json:"views"`
	Shares   int `json:"shares"`
}

// Pagination 分页信息
type Pagination struct {
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
	HasMore  bool  `json:"hasMore"`
}

// CreatePost 创建帖子
func (s *PostService) CreatePost(req *CreatePostRequest, authorId int64) (*CreatePostResponse, error) {
	// 验证分类是否存在
	category, err := s.categoryDao.GetByCode(req.Category)
	if err != nil {
		return nil, fmt.Errorf("分类不存在: %v", err)
	}

	// 处理标签和图片
	tagsJSON, _ := json.Marshal(req.Tags)
	imagesJSON, _ := json.Marshal(req.Images)

	// 生成摘要
	excerpt := req.Content
	if len(excerpt) > 200 {
		excerpt = excerpt[:200] + "..."
	}

	// 创建帖子
	post := &model.PostModel{
		Title:        req.Title,
		Content:      req.Content,
		Excerpt:      excerpt,
		AuthorId:     authorId,
		Category:     req.Category,
		CategoryName: category.Name,
		Tags:         string(tagsJSON),
		Images:       string(imagesJSON),
		IsPublic:     req.IsPublic,
	}

	err = s.postDao.Create(post)
	if err != nil {
		return nil, fmt.Errorf("创建帖子失败: %v", err)
	}

	// 更新分类帖子数量
	err = s.categoryDao.IncrementPostCount(req.Category)
	if err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("更新分类帖子数量失败: %v\n", err)
	}

	return &CreatePostResponse{
		PostId:    post.Id,
		CreatedAt: post.CreatedAt,
		URL:       fmt.Sprintf("https://api.example.com/posts/%d", post.Id),
	}, nil
}

// GetPostList 获取帖子列表
func (s *PostService) GetPostList(page, pageSize int, category, sort string, userId int64) (*PostListResponse, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 获取帖子列表
	posts, total, err := s.postDao.GetList(page, pageSize, category, sort)
	if err != nil {
		return nil, fmt.Errorf("获取帖子列表失败: %v", err)
	}

	// 获取用户点赞的帖子ID列表
	var likedPostIds []int64
	if userId != 0 {
		likedPostIds, err = s.userLikeDao.GetUserLikedPostIds(userId)
		if err != nil {
			// 记录错误但不影响主流程
			fmt.Printf("获取用户点赞列表失败: %v\n", err)
		}
	}

	// 构建响应数据
	postDetails := make([]*PostDetail, 0, len(posts))
	for _, post := range posts {
		// 获取作者信息
		author, err := s.userDao.GetById(post.AuthorId)
		if err != nil {
			// 如果获取作者信息失败，使用默认信息
			author = &model.UserModel{
				Id:         post.AuthorId,
				Nickname:   "未知用户",
				Avatar:     "",
				Bio:        "",
				Level:      1,
				IsVerified: false,
			}
		}

		// 解析标签和图片
		var tags []string
		var images []string
		json.Unmarshal([]byte(post.Tags), &tags)
		json.Unmarshal([]byte(post.Images), &images)

		// 检查是否点赞
		isLiked := false
		if userId != 0 {
			for _, likedId := range likedPostIds {
				if likedId == post.Id {
					isLiked = true
					break
				}
			}
		}

		postDetail := &PostDetail{
			Id:      post.Id,
			Title:   post.Title,
			Excerpt: post.Excerpt,
			Content: post.Content,
			Author: UserInfo{
				Id:         author.Id,
				Nickname:   author.Nickname,
				Avatar:     author.Avatar,
				Bio:        author.Bio,
				Level:      author.Level,
				IsVerified: author.IsVerified,
			},
			Category:     post.Category,
			CategoryName: post.CategoryName,
			Tags:         tags,
			Images:       images,
			Stats: PostStats{
				Likes:    post.Likes,
				Comments: post.Comments,
				Views:    post.Views,
				Shares:   post.Shares,
			},
			IsLiked:     isLiked,
			IsCollected: false, // TODO: 实现收藏功能
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		}
		postDetails = append(postDetails, postDetail)
	}

	// 计算分页信息
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	hasMore := page < totalPages

	return &PostListResponse{
		List: postDetails,
		Pagination: Pagination{
			Current:  page,
			PageSize: pageSize,
			Total:    total,
			HasMore:  hasMore,
		},
	}, nil
}
