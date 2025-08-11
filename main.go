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

	// 帖子相关接口
	http.HandleFunc("/api/posts", service.UserMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			postHandler.GetPostListHandler(w, r)
		case http.MethodPost:
			postHandler.CreatePostHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// 帖子详情和删除接口
	http.HandleFunc("/api/posts/", service.UserMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	// 分类相关接口
	http.HandleFunc("/api/categories", service.UserMiddleware(categoryHandler.GetCategoriesHandler))
	http.HandleFunc("/api/categories/publish", service.UserMiddleware(categoryHandler.GetPublishCategoriesHandler))
	http.HandleFunc("/api/topics/hot", service.UserMiddleware(categoryHandler.GetHotTopicsHandler))



	// 认证相关接口（不需要用户中间件）
	http.HandleFunc("/api/auth/", authHandler.HandleAuthRequests)

	// 用户相关接口
	http.HandleFunc("/api/user/", service.UserMiddleware(userHandler.HandleUserRequests))

	log.Fatal(http.ListenAndServe(":80", nil))
}
