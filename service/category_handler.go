package service

import (
	"encoding/json"
	"net/http"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService *CategoryService
}

// NewCategoryHandler 创建分类处理器实例
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: NewCategoryService(),
	}
}

// GetCategoriesHandler 获取分类列表处理器
func (h *CategoryHandler) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 调用服务
	result, err := h.categoryService.GetCategories()
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

// GetPublishCategoriesHandler 获取发布分类列表处理器
func (h *CategoryHandler) GetPublishCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 调用服务
	result, err := h.categoryService.GetPublishCategories()
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

// GetHotTopicsHandler 获取热门话题处理器
func (h *CategoryHandler) GetHotTopicsHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 获取用户ID（这里简化处理，实际应该从token中获取）
	userId := r.Header.Get("X-User-Id")

	// 调用服务
	result, err := h.categoryService.GetHotTopics(userId)
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