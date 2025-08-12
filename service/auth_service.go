package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// AuthService 认证服务
type AuthService struct {
	userDao dao.UserDao
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{
		userDao: dao.NewUserDao(),
	}
}

// RegisterUser 用户注册
func (s *AuthService) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// 获取微信小程序请求头
	openId := r.Header.Get("X-WX-OPENID")
	appId := r.Header.Get("X-WX-APPID")
	unionId := r.Header.Get("X-WX-UNIONID")
	fromOpenId := r.Header.Get("X-WX-FROM-OPENID")
	fromAppId := r.Header.Get("X-WX-FROM-APPID")
	fromUnionId := r.Header.Get("X-WX-FROM-UNIONID")

	// 优先使用 X-WX-OPENID，如果不存在则使用 X-WX-FROM-OPENID
	if openId == "" {
		openId = fromOpenId
	}
	if appId == "" {
		appId = fromAppId
	}
	if unionId == "" {
		unionId = fromUnionId
	}

	// 如果没有 openId，返回错误
	if openId == "" {
		http.Error(w, "Missing X-WX-OPENID header", http.StatusBadRequest)
		return
	}

	// 检查用户是否已存在
	existingUser, err := s.userDao.GetUserByOpenId(openId)
	if err == nil && existingUser != nil {
		// 用户已存在，返回用户信息
		response := map[string]interface{}{
			"code": 0,
			"msg":  "User already exists",
			"data": existingUser,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var registerData struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Bio      string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 生成随机用户名
	username := s.generateRandomUsername()

	// 创建新用户
	user := &model.UserModel{
		Username: username,
		Nickname: registerData.Nickname,
		Avatar:   registerData.Avatar,
		Bio:      registerData.Bio,
		OpenId:   openId,
		AppId:    appId,
		UnionId:  unionId,
		Level:    1,
		Password: "", // 微信小程序用户不需要密码
	}

	// 如果昵称为空，使用默认昵称
	if user.Nickname == "" {
		user.Nickname = "用户" + username
	}

	// 保存到数据库
	if err := s.userDao.CreateUser(user); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回创建成功的用户信息
	response := map[string]interface{}{
		"code": 0,
		"msg":  "User registered successfully",
		"data": user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LoginUser 用户登录（检查用户是否存在）
func (s *AuthService) LoginUser(w http.ResponseWriter, r *http.Request) {
	// 获取微信小程序请求头
	openId := r.Header.Get("X-WX-OPENID")
	fromOpenId := r.Header.Get("X-WX-FROM-OPENID")

	// 优先使用 X-WX-OPENID，如果不存在则使用 X-WX-FROM-OPENID
	if openId == "" {
		openId = fromOpenId
	}

	// 如果没有 openId，返回错误
	if openId == "" {
		http.Error(w, "Missing X-WX-OPENID header", http.StatusBadRequest)
		return
	}

	// 查询用户是否存在
	user, err := s.userDao.GetUserByOpenId(openId)
	if err != nil {
		// 用户不存在
		response := map[string]interface{}{
			"code": 1,
			"msg":  "User not found, please register first",
			"data": nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 用户存在，返回用户信息
	response := map[string]interface{}{
		"code": 0,
		"msg":  "Login successful",
		"data": user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckUserExists 检查用户是否存在
func (s *AuthService) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	// 获取微信小程序请求头
	openId := r.Header.Get("X-WX-OPENID")
	fromOpenId := r.Header.Get("X-WX-FROM-OPENID")

	// 优先使用 X-WX-OPENID，如果不存在则使用 X-WX-FROM-OPENID
	if openId == "" {
		openId = fromOpenId
	}

	// 如果没有 openId，返回错误
	if openId == "" {
		http.Error(w, "Missing X-WX-OPENID header", http.StatusBadRequest)
		return
	}

	// 查询用户是否存在
	user, _ := s.userDao.GetUserByOpenId(openId)

	if user == nil {
		http.Error(w, "user missing register first!", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"exists": user != nil,
			"user":   user,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateRandomUsername 生成随机用户名
func (s *AuthService) generateRandomUsername() string {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成6位随机数字
	randomNum := rand.Intn(900000) + 100000 // 100000-999999

	return fmt.Sprintf("user%d", randomNum)
}
