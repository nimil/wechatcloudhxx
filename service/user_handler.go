package service

import (
	"net/http"
	"strings"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: NewUserService(),
	}
}

// HandleUserRequests 处理用户相关请求
func (h *UserHandler) HandleUserRequests(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	switch {
	case path == "/api/user/profile":
		switch r.Method {
		case http.MethodGet:
			h.userService.GetCurrentUser(w, r)
		case http.MethodPut, http.MethodPatch:
			h.userService.UpdateUserProfile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case path == "/api/user/list":
		if r.Method == http.MethodGet {
			h.userService.GetUserList(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case strings.HasPrefix(path, "/api/user/"):
		// 处理其他用户相关请求
		if r.Method == http.MethodGet {
			h.userService.GetUserById(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}
