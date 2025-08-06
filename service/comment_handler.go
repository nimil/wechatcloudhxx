package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *CommentService
}

// NewCommentHandler 创建评论处理器实例
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: NewCommentService(),
	}
}

// CreateCommentHandler 创建评论处理器
func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 只允许POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL中提取帖子ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	postId := pathParts[3]

	// 解析请求体
	var req CreateCommentRequest
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
	result, err := h.commentService.CreateComment(postId, &req, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "评论成功",
		"data":    result,
	}

	json.NewEncoder(w).Encode(response)
}

// GetCommentListHandler 获取评论列表处理器
func (h *CommentHandler) GetCommentListHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 从URL中提取帖子ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	postId := pathParts[3]

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

	pageSize := 20
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 50 {
			pageSize = ps
		}
	}

	// 获取用户ID（这里简化处理，实际应该从token中获取）
	userId := r.Header.Get("X-User-Id")

	// 调用服务
	result, err := h.commentService.GetCommentList(postId, page, pageSize, userId)
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