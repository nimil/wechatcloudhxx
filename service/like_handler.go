package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// LikeHandler 点赞处理器
type LikeHandler struct {
	likeService *LikeService
}

// NewLikeHandler 创建点赞处理器实例
func NewLikeHandler() *LikeHandler {
	return &LikeHandler{
		likeService: NewLikeService(),
	}
}

// ToggleLikeHandler 切换点赞状态处理器
func (h *LikeHandler) ToggleLikeHandler(w http.ResponseWriter, r *http.Request) {
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
	postIdStr := pathParts[3]
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// 解析请求体
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 获取用户ID（这里简化处理，实际应该从token中获取）
	userIdStr := r.Header.Get("X-User-Id")
	if userIdStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 调用服务
	result, err := h.likeService.ToggleLike(postId, userId, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"code":    200,
		"message": "操作成功",
		"data":    result,
	}

	json.NewEncoder(w).Encode(response)
} 