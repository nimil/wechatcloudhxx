package service

import (
	"fmt"
	"math"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// CommentService 评论服务
type CommentService struct {
	commentDao dao.CommentDao
	userDao    dao.UserDao
	postDao    dao.PostDao
	securityService *ContentSecurityService
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{
		commentDao: dao.NewCommentDao(),
		userDao:    dao.NewUserDao(),
		postDao:    dao.NewPostDao(),
		securityService: NewContentSecurityService(),
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content"`
	ParentId int64  `json:"parentId"`
}

// CreateCommentResponse 创建评论响应
type CreateCommentResponse struct {
	CommentId int64     `json:"commentId"`
	CreatedAt time.Time `json:"createdAt"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	List       []*CommentDetail `json:"list"`
	Pagination Pagination       `json:"pagination"`
}

// CommentDetail 评论详情
type CommentDetail struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    UserInfo  `json:"author"`
	PostId    int64     `json:"postId"`
	ParentId  *int64    `json:"parentId"`
	Likes     int       `json:"likes"`
	IsLiked   bool      `json:"isLiked"`
	CreatedAt time.Time `json:"createdAt"`
	Replies   []*CommentDetail `json:"replies"`
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(postId int64, req *CreateCommentRequest, authorId int64, openid string) (*CreateCommentResponse, error) {
	// 验证帖子是否存在
	_, err := s.postDao.GetById(postId)
	if err != nil {
		return nil, fmt.Errorf("帖子不存在: %v", err)
	}

	// 内容安全校验
	if openid != "" && req.Content != "" {
		isSafe, err := s.securityService.IsContentSafe(openid, req.Content, SceneComment)
		if err != nil {
			return nil, fmt.Errorf("内容安全检测失败: %v", err)
		}
		if !isSafe {
			return nil, fmt.Errorf("评论内容包含违规信息，请修改后重试")
		}
	}

	// 创建评论
	comment := &model.CommentModel{
		Content:  req.Content,
		AuthorId: authorId,
		PostId:   postId,
		ParentId: nil,
	}

	// 如果有父评论ID，验证父评论是否存在
	if req.ParentId != 0 {
		_, err := s.commentDao.GetById(req.ParentId)
		if err != nil {
			return nil, fmt.Errorf("父评论不存在: %v", err)
		}
		comment.ParentId = &req.ParentId
	}

	err = s.commentDao.Create(comment)
	if err != nil {
		return nil, fmt.Errorf("创建评论失败: %v", err)
	}

	// 更新帖子评论数
	err = s.postDao.IncrementComments(postId)
	if err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("更新帖子评论数失败: %v\n", err)
	}

	return &CreateCommentResponse{
		CommentId: comment.Id,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// GetCommentList 获取评论列表
func (s *CommentService) GetCommentList(postId int64, page, pageSize int, userId int64) (*CommentListResponse, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	// 获取主评论列表
	comments, total, err := s.commentDao.GetByPostId(postId, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取评论列表失败: %v", err)
	}

	// 构建响应数据
	commentDetails := make([]*CommentDetail, 0, len(comments))
	for _, comment := range comments {
		// 获取作者信息
		author, err := s.userDao.GetById(comment.AuthorId)
		if err != nil {
			// 如果获取作者信息失败，使用默认信息
			author = &model.UserModel{
				Id:         comment.AuthorId,
				Nickname:   "未知用户",
				Avatar:     "",
				Bio:        "",
				Level:      1,
				IsVerified: false,
			}
		}

		commentDetail := &CommentDetail{
			Id:       comment.Id,
			Content:  comment.Content,
			Author: UserInfo{
				Id:         author.Id,
				Nickname:   author.Nickname,
				Avatar:     author.Avatar,
				Bio:        author.Bio,
				Level:      author.Level,
				IsVerified: author.IsVerified,
			},
			PostId:    comment.PostId,
			ParentId:  comment.ParentId,
			Likes:     comment.Likes,
			IsLiked:   false, // TODO: 实现评论点赞功能
			CreatedAt: comment.CreatedAt,
			Replies:   []*CommentDetail{}, // TODO: 实现回复功能
		}
		commentDetails = append(commentDetails, commentDetail)
	}

	// 计算分页信息
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	hasMore := page < totalPages

	return &CommentListResponse{
		List: commentDetails,
		Pagination: Pagination{
			Current:  page,
			PageSize: pageSize,
			Total:    total,
			HasMore:  hasMore,
		},
	}, nil
} 