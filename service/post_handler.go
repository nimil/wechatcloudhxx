package service

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	// 获取用户ID（这里简化处理，实际应该从token中获取）
	userId := r.Header.Get("X-User-Id")

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

	// 获取用户ID（这里简化处理，实际应该从token中获取）
	userId := r.Header.Get("X-User-Id")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 调用服务
	result, err := h.postService.CreatePost(&req, userId)
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
