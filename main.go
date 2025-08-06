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

	// 注册路由
	http.HandleFunc("/", service.IndexHandler)

	// 帖子相关接口
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			postHandler.GetPostListHandler(w, r)
		case http.MethodPost:
			postHandler.CreatePostHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 分类相关接口
	http.HandleFunc("/api/categories", categoryHandler.GetCategoriesHandler)
	http.HandleFunc("/api/categories/publish", categoryHandler.GetPublishCategoriesHandler)
	http.HandleFunc("/api/topics/hot", categoryHandler.GetHotTopicsHandler)

	// 评论相关接口
	http.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
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
			http.NotFound(w, r)
		}
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
