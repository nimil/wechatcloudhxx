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
	postDao         dao.PostDao
	userDao         dao.UserDao
	categoryDao     dao.CategoryDao
	userLikeDao     dao.UserLikeDao
	imageCheckDao   dao.ImageCheckDao
	securityService *ContentSecurityService
}

// NewPostService 创建帖子服务实例
func NewPostService() *PostService {
	return &PostService{
		postDao:         dao.NewPostDao(),
		userDao:         dao.NewUserDao(),
		categoryDao:     dao.NewCategoryDao(),
		userLikeDao:     dao.NewUserLikeDao(),
		imageCheckDao:   dao.NewImageCheckDao(),
		securityService: NewContentSecurityService(),
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
func (s *PostService) CreatePost(req *CreatePostRequest, authorId int64, openid string) (*CreatePostResponse, error) {
	// 验证分类是否存在
	category, err := s.categoryDao.GetByCode(req.Category)
	if err != nil {
		return nil, fmt.Errorf("分类不存在: %v", err)
	}

	// 内容安全校验
	if openid != "" {
		// 检查标题安全性（使用论坛场景）
		if req.Title != "" {
			isSafe, err := s.securityService.IsContentSafe(openid, req.Title, SceneForum)
			if err != nil {
				return nil, fmt.Errorf("标题安全检测失败: %v", err)
			}
			if !isSafe {
				return nil, fmt.Errorf("标题包含违规内容，请修改后重试")
			}
		}

		// 检查内容安全性（使用论坛场景）
		if req.Content != "" {
			isSafe, err := s.securityService.IsContentSafe(openid, req.Content, SceneForum)
			if err != nil {
				return nil, fmt.Errorf("内容安全检测失败: %v", err)
			}
			if !isSafe {
				return nil, fmt.Errorf("内容包含违规信息，请修改后重试")
			}
		}

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
		Title:            req.Title,
		Content:          req.Content,
		Excerpt:          excerpt,
		AuthorId:         authorId,
		Category:         req.Category,
		CategoryName:     category.Name,
		Tags:             string(tagsJSON),
		Images:           string(imagesJSON),
		ImageCheckStatus: 0, // 初始状态：待检测
		IsPublic:         req.IsPublic,
	}

	// 检查图片安全性（使用论坛场景）
	if len(req.Images) > 0 {
		// 先创建帖子，获取帖子ID
		err = s.postDao.Create(post)
		if err != nil {
			return nil, fmt.Errorf("创建帖子失败: %v", err)
		}

		// 设置帖子状态为检测中
		post.ImageCheckStatus = 1 // 检测中
		err = s.postDao.Update(post)
		if err != nil {
			return nil, fmt.Errorf("更新帖子状态失败: %v", err)
		}

		// 检测每张图片
		fmt.Printf("=== 开始检测帖子图片 - PostId: %d ===\n", post.Id)
		fmt.Printf("需要检测 %d 张图片\n", len(req.Images))
		
		for i, imageURL := range req.Images {
			if imageURL != "" {
				fmt.Printf("检测图片 %d: %s\n", i+1, imageURL)
				
				var result *MediaCheckResponse
				var err error
				
				// 判断是否为云存储文件ID，如果是则进行转换
				if s.securityService.cloudStorage.ValidateCloudID(imageURL) {
					fmt.Printf("检测到云存储文件ID，进行转换: %s\n", imageURL)
					
					// 使用云存储文件ID进行检测
					result, err = s.securityService.CheckCloudStorageImageSecurity(imageURL, openid, SceneForum)
					if err != nil {
						fmt.Printf("❌ 云存储图片%d检测失败: %v\n", i+1, err)
						return nil, fmt.Errorf("云存储图片安全检测失败: %v", err)
					}
				} else {
					// 直接使用URL进行检测
					result, err = s.securityService.CheckImageSecurity(imageURL, openid, SceneForum)
					if err != nil {
						fmt.Printf("❌ 图片%d检测失败: %v\n", i+1, err)
						return nil, fmt.Errorf("图片安全检测失败: %v", err)
					}
				}
				
				// 检查检测请求是否成功提交
				if !s.securityService.IsMediaCheckSuccess(result) {
					errorMsg := s.securityService.GetMediaCheckError(result)
					fmt.Printf("❌ 图片%d检测请求失败: %s\n", i+1, errorMsg)
					return nil, fmt.Errorf("图片安全检测请求失败: %s", errorMsg)
				}
				
				// 记录检测请求到数据库
				imageCheck := &model.ImageCheckModel{
					PostId:   post.Id,
					ImageURL: imageURL,
					TraceId:  result.TraceId,
					Status:   model.ImageCheckStatusChecking, // 检测中
				}
				
				err = s.imageCheckDao.Create(imageCheck)
				if err != nil {
					fmt.Printf("❌ 记录图片%d检测信息失败: %v\n", i+1, err)
					return nil, fmt.Errorf("记录图片检测信息失败: %v", err)
				}
				
				fmt.Printf("✅ 图片%d检测请求已提交，追踪ID: %s\n", i+1, result.TraceId)
			}
		}
		
		fmt.Printf("=== 图片检测请求提交完成 ===\n")

		// 返回帖子信息，但状态为检测中
		return &CreatePostResponse{
			PostId:    post.Id,
			CreatedAt: post.CreatedAt,
		}, nil
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
	}, nil
}

// GetPostDetail 获取帖子详情
func (s *PostService) GetPostDetail(postId int64, userId int64) (*PostDetail, error) {
	// 获取帖子信息
	post, err := s.postDao.GetById(postId)
	if err != nil {
		return nil, fmt.Errorf("帖子不存在: %v", err)
	}

	// 增加浏览量
	go func() {
		if err := s.postDao.IncrementViews(postId); err != nil {
			fmt.Printf("增加浏览量失败: %v\n", err)
		}
	}()

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
		likedPostIds, err := s.userLikeDao.GetUserLikedPostIds(userId)
		if err == nil {
			for _, likedId := range likedPostIds {
				if likedId == post.Id {
					isLiked = true
					break
				}
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

	return postDetail, nil
}

// SoftDeletePost 逻辑删除帖子
func (s *PostService) SoftDeletePost(postId int64, userId int64) error {
	// 获取帖子信息
	post, err := s.postDao.GetById(postId)
	if err != nil {
		return fmt.Errorf("帖子不存在: %v", err)
	}

	// 检查权限：只有作者可以删除自己的帖子
	if post.AuthorId != userId {
		return fmt.Errorf("无权限删除此帖子")
	}

	// 执行逻辑删除
	err = s.postDao.SoftDelete(postId)
	if err != nil {
		return fmt.Errorf("删除帖子失败: %v", err)
	}

	return nil
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

	// 获取帖子列表（只显示图片检测通过的帖子）
	posts, total, err := s.postDao.GetListWithImageCheck(page, pageSize, category, sort)
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
