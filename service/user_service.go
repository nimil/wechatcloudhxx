package service

import (
	"encoding/json"
	"net/http"
	"wxcloudrun-golang/db/dao"
)

// UserService 用户服务
type UserService struct {
	userDao dao.UserDao
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUserDao(),
	}
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userCtx := GetUserFromContext(r)
	if userCtx == nil || userCtx.User == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 返回用户信息
	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": userCtx.User,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userCtx := GetUserFromContext(r)
	if userCtx == nil || userCtx.User == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 解析请求体
	var updateData struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Bio      string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 更新用户信息
	if updateData.Nickname != "" {
		userCtx.User.Nickname = updateData.Nickname
	}
	if updateData.Avatar != "" {
		userCtx.User.Avatar = updateData.Avatar
	}
	if updateData.Bio != "" {
		userCtx.User.Bio = updateData.Bio
	}

	// 保存到数据库
	if err := s.userDao.UpdateUser(userCtx.User); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// 返回更新后的用户信息
	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": userCtx.User,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserById 根据ID获取用户信息
func (s *UserService) GetUserById(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中获取用户ID
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// 这里需要实现从字符串ID转换为int64的逻辑
	// 为了简化，这里先返回错误
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetUserList 获取用户列表（管理员功能）
func (s *UserService) GetUserList(w http.ResponseWriter, r *http.Request) {
	// 检查当前用户是否有管理员权限
	userCtx := GetUserFromContext(r)
	if userCtx == nil || userCtx.User == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 这里可以添加管理员权限检查
	// if userCtx.User.Level < 10 { // 假设10级以上为管理员
	//     http.Error(w, "Permission denied", http.StatusForbidden)
	//     return
	// }

	// 获取分页参数
	page := 1
	pageSize := 20

	users, total, err := s.userDao.GetUsersByPage(page, pageSize)
	if err != nil {
		http.Error(w, "Failed to get user list", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"users":     users,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
