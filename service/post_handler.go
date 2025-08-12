package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// PostHandler 帖子处理器
type PostHandler struct {
	postService *PostService
}

// NewPostHandler 创建帖子处理器实例
func NewPostHandler() *PostHandler {
	return &PostHandler{
		postService: NewPostService(),
	}
}

// GetPostListHandler 获取帖子列表处理器
func (h *PostHandler) GetPostListHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	category := r.URL.Query().Get("category")
	sort := r.URL.Query().Get("sort")

	// 解析分页参数
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 50 {
			pageSize = ps
		}
	}

	// 从用户上下文中获取用户ID
	userCtx := GetUserFromContext(r)
	var userId int64
	if userCtx != nil && userCtx.User != nil {
		userId = userCtx.User.Id
	}

	// 调用服务
	result, err := h.postService.GetPostList(page, pageSize, category, sort, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    result,
	}

	json.NewEncoder(w).Encode(response)
}

// CreatePostHandler 创建帖子处理器
func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 只允许POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 从用户上下文中获取用户ID
	userCtx := GetUserFromContext(r)
	if userCtx == nil || userCtx.User == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userId := userCtx.User.Id

	// 从请求头获取openid
	openid := r.Header.Get("x-wx-openid")

	// 调用服务
	result, err := h.postService.CreatePost(&req, userId, openid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "发布成功",
		"data":    result,
	}

	json.NewEncoder(w).Encode(response)
}

// GetPostDetailHandler 获取帖子详情处理器
func (h *PostHandler) GetPostDetailHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 从URL路径中提取帖子ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	postIdStr := pathParts[3]
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// 从用户上下文中获取用户ID
	userCtx := GetUserFromContext(r)
	var userId int64
	if userCtx != nil && userCtx.User != nil {
		userId = userCtx.User.Id
	}

	// 调用服务
	result, err := h.postService.GetPostDetail(postId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    result,
	}

	json.NewEncoder(w).Encode(response)
}

// DeletePostHandler 删除帖子处理器
func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 只允许DELETE请求
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中提取帖子ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	postIdStr := pathParts[3]
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// 从用户上下文中获取用户ID
	userCtx := GetUserFromContext(r)
	if userCtx == nil || userCtx.User == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userId := userCtx.User.Id

	// 调用服务
	err = h.postService.SoftDeletePost(postId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)
}

// GetMyPostsHandler 获取我的帖子处理器
func (h *PostHandler) GetMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 只允许GET请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// 解析分页参数
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 50 {
			pageSize = ps
		}
	}

	// 从用户上下文中获取用户ID
	userCtx := GetUserFromContext(r)
	if userCtx == nil || userCtx.User == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userId := userCtx.User.Id

	// 调用服务
	result, err := h.postService.GetUserPosts(userId, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    result,
	}

	json.NewEncoder(w).Encode(response)
}
