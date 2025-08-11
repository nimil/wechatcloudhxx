package service

import (
	"context"
	"net/http"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// UserContext 用户上下文
type UserContext struct {
	User   *model.UserModel
	OpenId string
	AppId  string
	UnionId string
	Env    string
	Source string
	IP     string
}

// UserMiddleware 用户认证中间件
func UserMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取微信小程序请求头
		openId := r.Header.Get("X-WX-OPENID")
		appId := r.Header.Get("X-WX-APPID")
		unionId := r.Header.Get("X-WX-UNIONID")
		fromOpenId := r.Header.Get("X-WX-FROM-OPENID")
		fromAppId := r.Header.Get("X-WX-FROM-APPID")
		fromUnionId := r.Header.Get("X-WX-FROM-UNIONID")
		env := r.Header.Get("X-WX-ENV")
		source := r.Header.Get("X-WX-SOURCE")
		ip := r.Header.Get("X-Original-Forwarded-For")

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
			http.Error(w, "Missing X-WX-OPENID header", http.StatusUnauthorized)
			return
		}

		// 创建用户DAO实例
		userDao := dao.NewUserDao()

			// 查询用户是否存在
	user, err := userDao.GetUserByOpenId(openId)
	if err != nil {
		// 用户不存在，返回错误
		http.Error(w, "User not found, please register first", http.StatusUnauthorized)
		return
	}

		// 创建用户上下文
		userCtx := &UserContext{
			User:    user,
			OpenId:  openId,
			AppId:   appId,
			UnionId: unionId,
			Env:     env,
			Source:  source,
			IP:      ip,
		}

		// 将用户上下文存储到请求上下文中
		ctx := context.WithValue(r.Context(), "user", userCtx)
		r = r.WithContext(ctx)

		// 调用下一个处理器
		next(w, r)
	}
}

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(r *http.Request) *UserContext {
	if userCtx, ok := r.Context().Value("user").(*UserContext); ok {
		return userCtx
	}
	return nil
}



// RequireAuth 要求认证的中间件
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := GetUserFromContext(r)
		if userCtx == nil || userCtx.User == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
