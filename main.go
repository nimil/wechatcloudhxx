package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// 创建处理器实例
	postHandler := service.NewPostHandler()
	categoryHandler := service.NewCategoryHandler()
	commentHandler := service.NewCommentHandler()
	likeHandler := service.NewLikeHandler()
	userHandler := service.NewUserHandler()
	authHandler := service.NewAuthHandler()


	// 统一的帖子路由处理器
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// 精确匹配 /api/posts
		if path == "/api/posts" {
			switch r.Method {
			case http.MethodGet:
				// GET 请求不需要认证，直接处理
				postHandler.GetPostListHandler(w, r)
			case http.MethodPost:
				// POST 请求需要认证
				service.UserMiddleware(postHandler.CreatePostHandler)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		
		// 处理 /api/posts/ 开头的路径（需要认证）
		service.UserMiddleware(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if strings.HasSuffix(path, "/comments") {
				switch r.Method {
				case http.MethodGet:
					commentHandler.GetCommentListHandler(w, r)
				case http.MethodPost:
					commentHandler.CreateCommentHandler(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			} else if strings.HasSuffix(path, "/like") {
				likeHandler.ToggleLikeHandler(w, r)
			} else {
				// 处理帖子详情和删除
				switch r.Method {
				case http.MethodGet:
					postHandler.GetPostDetailHandler(w, r)
				case http.MethodDelete:
					postHandler.DeletePostHandler(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}
		})(w, r)
	})

	// 分类相关接口
	http.HandleFunc("/api/categories", service.UserMiddleware(categoryHandler.GetCategoriesHandler))
	http.HandleFunc("/api/categories/publish", service.UserMiddleware(categoryHandler.GetPublishCategoriesHandler))
	http.HandleFunc("/api/topics/hot", categoryHandler.GetHotTopicsHandler)



	// 认证相关接口（不需要用户中间件）
	http.HandleFunc("/api/auth/", authHandler.HandleAuthRequests)

	// 用户相关接口
	http.HandleFunc("/api/user/", service.UserMiddleware(userHandler.HandleUserRequests))

	// 微信回调接口（不需要用户中间件）
	wechatCallbackHandler := service.NewWechatCallbackHandler()
	http.HandleFunc("/api/wechat/callback", wechatCallbackHandler.HandleMediaCheckCallback)

	log.Fatal(http.ListenAndServe(":80", nil))
}
