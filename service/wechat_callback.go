package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

// WechatMediaCheckCallback 微信媒体检测回调数据结构
type WechatMediaCheckCallback struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Event        string `json:"Event"`
	Appid        string `json:"appid"`
	TraceId      string `json:"trace_id"`
	Version      int    `json:"version"`
	Detail       []struct {
		Strategy string  `json:"strategy"`
		Errcode  int     `json:"errcode"`
		Suggest  string  `json:"suggest"`
		Label    int     `json:"label"`
		Prob     float64 `json:"prob"`
	} `json:"detail"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Result  struct {
		Suggest string `json:"suggest"`
		Label   int    `json:"label"`
	} `json:"result"`
}

// WechatCallbackHandler 微信回调处理器
type WechatCallbackHandler struct {
	imageCheckDao dao.ImageCheckDao
	postDao       dao.PostDao
}

// NewWechatCallbackHandler 创建微信回调处理器
func NewWechatCallbackHandler() *WechatCallbackHandler {
	return &WechatCallbackHandler{
		imageCheckDao: dao.NewImageCheckDao(),
		postDao:       dao.NewPostDao(),
	}
}

// HandleMediaCheckCallback 处理媒体检测回调
func (h *WechatCallbackHandler) HandleMediaCheckCallback(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 只允许POST请求
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("读取回调请求体失败: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// 打印回调请求详情
	fmt.Printf("=== 微信回调请求详情 ===\n")
	fmt.Printf("请求方法: %s\n", r.Method)
	fmt.Printf("请求URL: %s\n", r.URL.String())
	fmt.Printf("请求头:\n")
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	fmt.Printf("请求体: %s\n", string(body))
	fmt.Printf("=== 微信回调请求详情结束 ===\n")

	// 首先尝试解析为验证请求格式
	var verifyRequest struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(body, &verifyRequest); err == nil && verifyRequest.Action == "CheckContainerPath" {
		// 这是验证请求，直接返回成功
		log.Printf("收到验证请求: %s", verifyRequest.Action)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
		return
	}

	// 解析回调数据
	var callback WechatMediaCheckCallback
	if err := json.Unmarshal(body, &callback); err != nil {
		log.Printf("解析回调数据失败: %v", err)
		http.Error(w, "Invalid callback data", http.StatusBadRequest)
		return
	}

	// 验证回调类型
	if callback.Event != "wxa_media_check" {
		log.Printf("未知的回调事件类型: %s", callback.Event)
		http.Error(w, "Unknown event type", http.StatusBadRequest)
		return
	}

	// 处理媒体检测结果
	err = h.processMediaCheckResult(&callback)
	if err != nil {
		log.Printf("处理媒体检测结果失败: %v", err)
		http.Error(w, "Failed to process media check result", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// processMediaCheckResult 处理媒体检测结果
func (h *WechatCallbackHandler) processMediaCheckResult(callback *WechatMediaCheckCallback) error {
	fmt.Printf("=== 处理媒体检测回调 ===\n")
	fmt.Printf("回调数据详情:\n")
	fmt.Printf("  - ToUserName: %s\n", callback.ToUserName)
	fmt.Printf("  - FromUserName: %s\n", callback.FromUserName)
	fmt.Printf("  - CreateTime: %d\n", callback.CreateTime)
	fmt.Printf("  - MsgType: %s\n", callback.MsgType)
	fmt.Printf("  - Event: %s\n", callback.Event)
	fmt.Printf("  - Appid: %s\n", callback.Appid)
	fmt.Printf("  - TraceId: %s\n", callback.TraceId)
	fmt.Printf("  - Version: %d\n", callback.Version)
	fmt.Printf("  - Errcode: %d\n", callback.Errcode)
	fmt.Printf("  - Errmsg: %s\n", callback.Errmsg)
	fmt.Printf("  - Result.Suggest: %s\n", callback.Result.Suggest)
	fmt.Printf("  - Result.Label: %d\n", callback.Result.Label)
	
	if len(callback.Detail) > 0 {
		fmt.Printf("  - Detail数组:\n")
		for i, detail := range callback.Detail {
			fmt.Printf("    [%d] Strategy: %s, Errcode: %d, Suggest: %s, Label: %d, Prob: %.2f\n",
				i, detail.Strategy, detail.Errcode, detail.Suggest, detail.Label, detail.Prob)
		}
	} else {
		fmt.Printf("  - Detail数组: 空\n")
	}
	fmt.Printf("=== 回调数据详情结束 ===\n")

	// 根据trace_id查找对应的检测记录
	fmt.Printf("查找检测记录 - TraceId: %s\n", callback.TraceId)
	imageCheck, err := h.imageCheckDao.GetByTraceId(callback.TraceId)
	if err != nil {
		fmt.Printf("❌ 未找到对应的检测记录: %v\n", err)
		return fmt.Errorf("未找到对应的检测记录: %v", err)
	}
	fmt.Printf("✅ 找到检测记录 - PostId: %d, ImageURL: %s, Status: %d\n", 
		imageCheck.PostId, imageCheck.ImageURL, imageCheck.Status)

	// 确定检测状态和建议
	var status int
	var suggest string
	var label int
	var prob float64
	var strategy string

	if callback.Errcode == 0 {
		// 检测成功
		if len(callback.Detail) > 0 {
			detail := callback.Detail[0]
			suggest = detail.Suggest
			label = detail.Label
			prob = detail.Prob
			strategy = detail.Strategy
		} else {
			suggest = callback.Result.Suggest
			label = callback.Result.Label
		}

		// 根据建议确定状态
		switch suggest {
		case "pass":
			status = model.ImageCheckStatusPassed
		case "review":
			status = model.ImageCheckStatusFailed
		case "risky":
			status = model.ImageCheckStatusFailed
		default:
			status = model.ImageCheckStatusFailed
		}
	} else {
		// 检测失败
		status = model.ImageCheckStatusFailed
		suggest = "failed"
		label = 0
		prob = 0
		strategy = ""
	}

	// 更新检测记录
	fmt.Printf("更新检测记录:\n")
	fmt.Printf("  - TraceId: %s\n", callback.TraceId)
	fmt.Printf("  - Status: %d\n", status)
	fmt.Printf("  - Suggest: %s\n", suggest)
	fmt.Printf("  - Label: %d\n", label)
	fmt.Printf("  - Prob: %.2f\n", prob)
	fmt.Printf("  - Strategy: %s\n", strategy)
	fmt.Printf("  - Errcode: %d\n", callback.Errcode)
	fmt.Printf("  - Errmsg: %s\n", callback.Errmsg)
	
	err = h.imageCheckDao.UpdateStatus(
		callback.TraceId,
		status,
		suggest,
		label,
		prob,
		strategy,
		callback.Errcode,
		callback.Errmsg,
	)
	if err != nil {
		fmt.Printf("❌ 更新检测记录失败: %v\n", err)
		return fmt.Errorf("更新检测记录失败: %v", err)
	}
	fmt.Printf("✅ 检测记录更新成功\n")

	// 检查帖子的所有图片是否都检测完成
	err = h.checkPostImageCheckStatus(imageCheck.PostId)
	if err != nil {
		return fmt.Errorf("检查帖子图片检测状态失败: %v", err)
	}

	log.Printf("媒体检测结果处理完成 - TraceId: %s, Status: %d, Suggest: %s",
		callback.TraceId, status, suggest)

	return nil
}

// checkPostImageCheckStatus 检查帖子的所有图片检测状态
func (h *WechatCallbackHandler) checkPostImageCheckStatus(postId int64) error {
	fmt.Printf("=== 检查帖子图片检测状态 - PostId: %d ===\n", postId)
	
	// 获取帖子的所有图片检测记录
	imageChecks, err := h.imageCheckDao.GetByPostId(postId)
	if err != nil {
		fmt.Printf("❌ 获取帖子图片检测记录失败: %v\n", err)
		return fmt.Errorf("获取帖子图片检测记录失败: %v", err)
	}

	fmt.Printf("帖子共有 %d 张图片需要检测\n", len(imageChecks))
	
	if len(imageChecks) == 0 {
		fmt.Printf("没有图片，无需处理\n")
		return nil // 没有图片，无需处理
	}

	// 检查是否所有图片都检测完成
	allCompleted := true
	allPassed := true

	fmt.Printf("检查每张图片的检测状态:\n")
	for i, check := range imageChecks {
		statusText := "未知"
		switch check.Status {
		case model.ImageCheckStatusPending:
			statusText = "待检测"
		case model.ImageCheckStatusChecking:
			statusText = "检测中"
		case model.ImageCheckStatusPassed:
			statusText = "检测通过"
		case model.ImageCheckStatusFailed:
			statusText = "检测失败"
		}
		
		fmt.Printf("  [%d] 图片: %s, 状态: %d (%s), 建议: %s\n", 
			i+1, check.ImageURL, check.Status, statusText, check.Suggest)
		
		if check.Status == model.ImageCheckStatusPending || check.Status == model.ImageCheckStatusChecking {
			allCompleted = false
		}
		if check.Status == model.ImageCheckStatusFailed {
			allPassed = false
		}
	}

	fmt.Printf("检测完成状态: %t, 全部通过: %t\n", allCompleted, allPassed)

	if !allCompleted {
		fmt.Printf("还有图片在检测中，等待下次回调\n")
		return nil // 还有图片在检测中
	}

	// 所有图片检测完成，更新帖子状态
	var postStatus int
	var statusText string
	if allPassed {
		postStatus = model.ImageCheckStatusPassed // 所有图片检测通过
		statusText = "检测通过"
	} else {
		postStatus = model.ImageCheckStatusFailed // 有图片检测失败
		statusText = "检测失败"
	}

	fmt.Printf("更新帖子状态: %d (%s)\n", postStatus, statusText)

	// 更新帖子状态
	err = h.postDao.UpdateImageCheckStatus(postId, postStatus)
	if err != nil {
		fmt.Printf("❌ 更新帖子图片检测状态失败: %v\n", err)
		return fmt.Errorf("更新帖子图片检测状态失败: %v", err)
	}

	fmt.Printf("✅ 帖子图片检测完成 - PostId: %d, Status: %d (%s)\n", postId, postStatus, statusText)
	fmt.Printf("=== 检查帖子图片检测状态结束 ===\n")

	return nil
}
