package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// WechatCloudStorageService 微信云存储服务
type WechatCloudStorageService struct {
	client *http.Client
}

// NewWechatCloudStorageService 创建微信云存储服务实例
func NewWechatCloudStorageService() *WechatCloudStorageService {
	return &WechatCloudStorageService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// BatchDownloadFileRequest 批量下载文件请求结构
type BatchDownloadFileRequest struct {
	Env      string `json:"env"`
	FileList []struct {
		FileID  string `json:"fileid"`
		MaxAge  int    `json:"max_age"`
	} `json:"file_list"`
}

// FileDownloadInfo 文件下载信息
type FileDownloadInfo struct {
	FileID       string `json:"fileid"`
	DownloadURL  string `json:"download_url"`
	Status       int    `json:"status"`
	Errmsg       string `json:"errmsg"`
}

// BatchDownloadFileResponse 批量下载文件响应结构
type BatchDownloadFileResponse struct {
	Errcode  int                `json:"errcode"`
	Errmsg   string             `json:"errmsg"`
	FileList []FileDownloadInfo `json:"file_list"`
}

// GetFileDownloadURL 获取单个文件的下载URL
func (s *WechatCloudStorageService) GetFileDownloadURL(cloudID string) (string, error) {
	// 获取环境ID
	envID := os.Getenv("ENV_ID")
	if envID == "" {
		envID = "werun-id" // 默认环境ID，本地调试时使用
	}

	// 构建请求数据
	requestData := BatchDownloadFileRequest{
		Env: envID,
		FileList: []struct {
			FileID  string `json:"fileid"`
			MaxAge  int    `json:"max_age"`
		}{
			{
				FileID: cloudID,
				MaxAge: 86400, // 24小时有效期
			},
		},
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "http://api.weixin.qq.com/tcb/batchdownloadfile", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 打印请求详情
	fmt.Printf("=== 微信云存储下载URL请求详情 ===\n")
	fmt.Printf("请求URL: %s\n", req.URL.String())
	fmt.Printf("请求方法: %s\n", req.Method)
	fmt.Printf("环境ID: %s\n", envID)
	fmt.Printf("云存储文件ID: %s\n", cloudID)
	fmt.Printf("请求体JSON: %s\n", string(jsonData))
	fmt.Printf("=== 微信云存储下载URL请求详情结束 ===\n")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 打印响应详情
	fmt.Printf("=== 微信云存储下载URL响应详情 ===\n")
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应体: %s\n", string(body))
	fmt.Printf("=== 微信云存储下载URL响应详情结束 ===\n")

	// 解析响应
	var response BatchDownloadFileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查错误码
	if response.Errcode != 0 {
		return "", fmt.Errorf("微信云存储API错误: %s (错误码: %d)", response.Errmsg, response.Errcode)
	}

	// 检查文件列表
	if len(response.FileList) == 0 {
		return "", fmt.Errorf("未找到文件下载信息")
	}

	fileInfo := response.FileList[0]
	if fileInfo.Status != 0 {
		return "", fmt.Errorf("文件下载状态错误: %s", fileInfo.Errmsg)
	}

	if fileInfo.DownloadURL == "" {
		return "", fmt.Errorf("未获取到下载URL")
	}

	return fileInfo.DownloadURL, nil
}

// GetMultipleFileDownloadURLs 获取多个文件的下载URL
func (s *WechatCloudStorageService) GetMultipleFileDownloadURLs(cloudIDs []string) (map[string]string, error) {
	if len(cloudIDs) == 0 {
		return nil, fmt.Errorf("文件ID列表不能为空")
	}

	// 获取环境ID
	envID := os.Getenv("ENV_ID")
	if envID == "" {
		envID = "werun-id" // 默认环境ID，本地调试时使用
	}

	// 构建文件列表
	fileList := make([]struct {
		FileID  string `json:"fileid"`
		MaxAge  int    `json:"max_age"`
	}, len(cloudIDs))

	for i, cloudID := range cloudIDs {
		fileList[i] = struct {
			FileID  string `json:"fileid"`
			MaxAge  int    `json:"max_age"`
		}{
			FileID: cloudID,
			MaxAge: 86400, // 24小时有效期
		}
	}

	// 构建请求数据
	requestData := BatchDownloadFileRequest{
		Env:      envID,
		FileList: fileList,
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "http://api.weixin.qq.com/tcb/batchdownloadfile", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 打印请求详情
	fmt.Printf("=== 微信云存储批量下载URL请求详情 ===\n")
	fmt.Printf("请求URL: %s\n", req.URL.String())
	fmt.Printf("请求方法: %s\n", req.Method)
	fmt.Printf("环境ID: %s\n", envID)
	fmt.Printf("文件数量: %d\n", len(cloudIDs))
	fmt.Printf("请求体JSON: %s\n", string(jsonData))
	fmt.Printf("=== 微信云存储批量下载URL请求详情结束 ===\n")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 打印响应详情
	fmt.Printf("=== 微信云存储批量下载URL响应详情 ===\n")
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应体: %s\n", string(body))
	fmt.Printf("=== 微信云存储批量下载URL响应详情结束 ===\n")

	// 解析响应
	var response BatchDownloadFileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查错误码
	if response.Errcode != 0 {
		return nil, fmt.Errorf("微信云存储API错误: %s (错误码: %d)", response.Errmsg, response.Errcode)
	}

	// 构建结果映射
	result := make(map[string]string)
	for _, fileInfo := range response.FileList {
		if fileInfo.Status == 0 && fileInfo.DownloadURL != "" {
			result[fileInfo.FileID] = fileInfo.DownloadURL
		} else {
			fmt.Printf("警告: 文件 %s 获取下载URL失败: %s\n", fileInfo.FileID, fileInfo.Errmsg)
		}
	}

	return result, nil
}

// ValidateCloudID 验证云存储文件ID格式
func (s *WechatCloudStorageService) ValidateCloudID(cloudID string) bool {
	// 基本的云存储文件ID格式验证
	// 通常格式为: cloud://env-id.app-id-env-id/filepath
	if cloudID == "" {
		return false
	}
	
	// 检查是否以 "cloud://" 开头
	if len(cloudID) < 8 || cloudID[:8] != "cloud://" {
		return false
	}
	
	return true
}

// GetEnvironmentID 获取当前环境ID
func (s *WechatCloudStorageService) GetEnvironmentID() string {
	envID := os.Getenv("ENV_ID")
	if envID == "" {
		envID = "werun-id" // 默认环境ID
	}
	return envID
}
