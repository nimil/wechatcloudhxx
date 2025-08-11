package service

import (
	"net/http"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: NewAuthService(),
	}
}

// HandleAuthRequests 处理认证相关请求
func (h *AuthHandler) HandleAuthRequests(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	switch {
	case path == "/api/auth/register":
		if r.Method == http.MethodPost {
			h.authService.RegisterUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case path == "/api/auth/login":
		if r.Method == http.MethodPost {
			h.authService.LoginUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case path == "/api/auth/check":
		if r.Method == http.MethodGet {
			h.authService.CheckUserExists(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}
