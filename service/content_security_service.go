package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 内容安全检测场景值常量
const (
	SceneProfile = 1 // 资料
	SceneComment = 2 // 评论
	SceneForum   = 3 // 论坛
	SceneSocial  = 4 // 社交日志
)

// 内容安全检测标签值常量
const (
	LabelNormal     = 100   // 正常/垃圾信息
	LabelPorn       = 20001 // 色情
	LabelAbuse      = 20002 // 辱骂
	LabelPolitics   = 20003 // 政治
	LabelAd         = 20006 // 广告
	LabelTerrorism  = 20008 // 违法犯罪
	LabelOther      = 20012 // 其他
)

// 内容安全检测建议值常量
const (
	SuggestPass   = "pass"   // 通过
	SuggestReview = "review" // 需要人工审核
	SuggestRisky  = "risky"  // 有风险
)

// ContentSecurityService 内容安全校验服务
type ContentSecurityService struct {
	client *http.Client
}

// NewContentSecurityService 创建内容安全校验服务实例
func NewContentSecurityService() *ContentSecurityService {
	return &ContentSecurityService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// MsgSecCheckRequest 内容安全检测请求
type MsgSecCheckRequest struct {
	Openid  string `json:"openid"`
	Version int    `json:"version"`
	Scene   int    `json:"scene"`
	Content string `json:"content"`
}

// MsgSecCheckDetail 内容安全检测详细结果
type MsgSecCheckDetail struct {
	Strategy string  `json:"strategy"`
	Errcode  int     `json:"errcode"`
	Suggest  string  `json:"suggest"`
	Label    int     `json:"label"`
	Prob     float64 `json:"prob,omitempty"`
	Level    int     `json:"level,omitempty"`
	Keyword  string  `json:"keyword,omitempty"`
}

// MsgSecCheckResponse 内容安全检测响应
type MsgSecCheckResponse struct {
	Errcode int                  `json:"errcode"`
	Errmsg  string               `json:"errmsg"`
	TraceId string               `json:"trace_id,omitempty"`
	Result  struct {
		Suggest string `json:"suggest"`
		Label   int    `json:"label"`
	} `json:"result,omitempty"`
	Detail []MsgSecCheckDetail `json:"detail,omitempty"`
}

// CheckContentSecurity 检查内容安全性
func (s *ContentSecurityService) CheckContentSecurity(openid, content string, scene int) (*MsgSecCheckResponse, error) {
	// 构建请求数据
	requestData := MsgSecCheckRequest{
		Openid:  openid,
		Version: 2,
		Scene:   scene, // 场景值：1 资料；2 评论；3 论坛；4 社交日志
		Content: content,
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	//
	req, err := http.NewRequest("POST", "http://api.weixin.qq.com/wxa/msg_sec_check", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

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

	// 解析响应
	var response MsgSecCheckResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

// IsContentSafe 判断内容是否安全
func (s *ContentSecurityService) IsContentSafe(openid, content string, scene int) (bool, error) {
	response, err := s.CheckContentSecurity(openid, content, scene)
	if err != nil {
		return false, err
	}

	// 检查错误码
	if response.Errcode != 0 {
		return false, fmt.Errorf("内容安全检测失败: %s", response.Errmsg)
	}

	// 只判断suggest字段，只有pass时通过
	isSafe := response.Result.Suggest == "pass"
	
	// 打印检测结果日志
	fmt.Printf("[内容安全检测] OpenID: %s, 场景: %d, 建议: %s, 标签: %d, 是否通过: %t\n", 
		openid, scene, response.Result.Suggest, response.Result.Label, isSafe)
	
	return isSafe, nil
}

// GetContentSecurityResult 获取内容安全检测详细结果
func (s *ContentSecurityService) GetContentSecurityResult(openid, content string, scene int) (*MsgSecCheckResponse, error) {
	return s.CheckContentSecurity(openid, content, scene)
}

// GetContentSecurityDetail 获取内容安全检测的详细分析
func (s *ContentSecurityService) GetContentSecurityDetail(openid, content string, scene int) (*ContentSecurityDetail, error) {
	response, err := s.CheckContentSecurity(openid, content, scene)
	if err != nil {
		return nil, err
	}

	// 只判断suggest字段，只有pass时通过
	isSafe := response.Result.Suggest == "pass"
	
	// 打印检测结果日志
	fmt.Printf("[内容安全检测详情] OpenID: %s, 场景: %d, 建议: %s, 标签: %d, 是否通过: %t, 追踪ID: %s\n", 
		openid, scene, response.Result.Suggest, response.Result.Label, isSafe, response.TraceId)

	detail := &ContentSecurityDetail{
		IsSafe:    isSafe,
		Suggest:   response.Result.Suggest,
		Label:     response.Result.Label,
		TraceId:   response.TraceId,
		Errcode:   response.Errcode,
		Errmsg:    response.Errmsg,
		Details:   response.Detail,
		RiskLevel: s.calculateRiskLevel(response),
		Keywords:  s.extractKeywords(response.Detail),
	}

	return detail, nil
}

// ContentSecurityDetail 内容安全检测详细分析结果
type ContentSecurityDetail struct {
	IsSafe    bool                 `json:"isSafe"`
	Suggest   string               `json:"suggest"`
	Label     int                  `json:"label"`
	TraceId   string               `json:"traceId"`
	Errcode   int                  `json:"errcode"`
	Errmsg    string               `json:"errmsg"`
	Details   []MsgSecCheckDetail  `json:"details"`
	RiskLevel string               `json:"riskLevel"`
	Keywords  []string             `json:"keywords"`
}

// calculateRiskLevel 计算风险等级
func (s *ContentSecurityService) calculateRiskLevel(response *MsgSecCheckResponse) string {
	if response.Result.Suggest == "pass" {
		return "safe"
	}

	// 分析详细结果中的最高风险等级
	maxLevel := 0
	for _, detail := range response.Detail {
		if detail.Level > maxLevel {
			maxLevel = detail.Level
		}
	}

	switch {
	case maxLevel >= 90:
		return "high"
	case maxLevel >= 60:
		return "medium"
	default:
		return "low"
	}
}

// extractKeywords 提取命中的关键词
func (s *ContentSecurityService) extractKeywords(details []MsgSecCheckDetail) []string {
	var keywords []string
	for _, detail := range details {
		if detail.Keyword != "" {
			keywords = append(keywords, detail.Keyword)
		}
	}
	return keywords
}

// GetLabelDescription 获取标签描述
func (s *ContentSecurityService) GetLabelDescription(label int) string {
	switch label {
	case LabelNormal:
		return "正常/垃圾信息"
	case LabelPorn:
		return "色情"
	case LabelAbuse:
		return "辱骂"
	case LabelPolitics:
		return "政治"
	case LabelAd:
		return "广告"
	case LabelTerrorism:
		return "违法犯罪"
	case LabelOther:
		return "其他"
	default:
		return "未知"
	}
}

// GetSuggestDescription 获取建议描述
func (s *ContentSecurityService) GetSuggestDescription(suggest string) string {
	switch suggest {
	case SuggestPass:
		return "通过"
	case SuggestReview:
		return "需要人工审核"
	case SuggestRisky:
		return "有风险"
	default:
		return "未知"
	}
}
