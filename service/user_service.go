package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

// UserRequest 用户请求结构
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse 用户响应结构
type UserResponse struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	OpenId   string `json:"openid,omitempty"`
	UnionId  string `json:"unionid,omitempty"`
	AppId    string `json:"appid,omitempty"`
}

// UserListResponse 用户列表响应结构
type UserListResponse struct {
	Users    []*UserResponse `json:"users"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Pages    int             `json:"pages"`
}

// UserHandler 用户接口
func UserHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodPost {
		user, err := createUser(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = user
		}
	} else if r.Method == http.MethodGet {
		// 分页查询用户列表
		userList, err := getUsersByPage(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = userList
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// createUser 创建用户
func createUser(r *http.Request) (*UserResponse, error) {
	// 解析请求参数
	var userReq UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		return nil, fmt.Errorf("解析请求参数失败: %v", err)
	}
	defer r.Body.Close()

	// 验证参数
	if userReq.Username == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}
	if userReq.Password == "" {
		return nil, fmt.Errorf("密码不能为空")
	}

	// 获取微信用户信息
	openId := r.Header.Get("X-WX-OPENID")
	unionId := r.Header.Get("X-WX-UNIONID")
	appId := r.Header.Get("X-WX-APPID")

	// 如果资源复用，使用FROM字段
	if openId == "" {
		openId = r.Header.Get("X-WX-FROM-OPENID")
	}
	if unionId == "" {
		unionId = r.Header.Get("X-WX-FROM-UNIONID")
	}
	if appId == "" {
		appId = r.Header.Get("X-WX-FROM-APPID")
	}

	// 检查用户是否已存在（优先检查微信信息）
	if openId != "" {
		existingUser, err := dao.UserImp.GetUserByOpenId(openId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("查询用户失败: %v", err)
		}
		if existingUser != nil && existingUser.Id > 0 {
			return nil, fmt.Errorf("该微信用户已存在")
		}
	}

	if unionId != "" {
		existingUser, err := dao.UserImp.GetUserByUnionId(unionId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("查询用户失败: %v", err)
		}
		if existingUser != nil && existingUser.Id > 0 {
			return nil, fmt.Errorf("该微信用户已存在")
		}
	}

	// 检查用户名是否已存在
	existingUser, err := dao.UserImp.GetUserByUsername(userReq.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	if existingUser != nil && existingUser.Id > 0 {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 创建新用户
	now := time.Now()
	user := &model.UserModel{
		Username:  userReq.Username,
		Password:  userReq.Password, // 注意：实际项目中应该对密码进行加密
		OpenId:    openId,
		UnionId:   unionId,
		AppId:     appId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = dao.UserImp.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 返回用户信息（不包含密码）
	response := &UserResponse{
		Id:       user.Id,
		Username: user.Username,
		OpenId:   user.OpenId,
		UnionId:  user.UnionId,
		AppId:    user.AppId,
	}

	return response, nil
}

// getUsersByPage 分页查询用户列表
func getUsersByPage(r *http.Request) (*UserListResponse, error) {
	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// 设置默认值
	page := 1
	pageSize := 10

	// 解析页码
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			return nil, fmt.Errorf("页码参数无效")
		}
	}

	// 解析每页大小
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		} else {
			return nil, fmt.Errorf("每页大小参数无效，范围1-100")
		}
	}

	// 查询用户列表
	users, total, err := dao.UserImp.GetUsersByPage(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %v", err)
	}

	// 转换为响应格式
	userResponses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, &UserResponse{
			Id:       user.Id,
			Username: user.Username,
			OpenId:   user.OpenId,
			UnionId:  user.UnionId,
			AppId:    user.AppId,
		})
	}

	// 计算总页数
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := &UserListResponse{
		Users:    userResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}

	return response, nil
}
