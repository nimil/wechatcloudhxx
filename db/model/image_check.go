package model

import "time"

// ImageCheckModel 图片检测记录模型
type ImageCheckModel struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostId      int64     `gorm:"column:post_id;not null;index" json:"postId"`      // 关联的帖子ID
	ImageURL    string    `gorm:"column:image_url;type:varchar(500);not null" json:"imageUrl"` // 图片URL
	TraceId     string    `gorm:"column:trace_id;type:varchar(100);not null;index" json:"traceId"` // 微信检测追踪ID
	Status      int       `gorm:"column:status;default:0" json:"status"` // 检测状态：0-待检测 1-检测中 2-检测通过 3-检测失败
	Suggest     string    `gorm:"column:suggest;type:varchar(20)" json:"suggest"` // 检测建议：pass/review/risky
	Label       int       `gorm:"column:label;default:0" json:"label"` // 检测标签
	Prob        float64   `gorm:"column:prob;type:decimal(5,2)" json:"prob"` // 置信度
	Strategy    string    `gorm:"column:strategy;type:varchar(50)" json:"strategy"` // 检测策略
	Errcode     int       `gorm:"column:errcode;default:0" json:"errcode"` // 错误码
	Errmsg      string    `gorm:"column:errmsg;type:varchar(200)" json:"errmsg"` // 错误信息
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (ImageCheckModel) TableName() string {
	return "image_checks"
}

// 图片检测状态常量
const (
	ImageCheckStatusPending = 0 // 待检测
	ImageCheckStatusChecking = 1 // 检测中
	ImageCheckStatusPassed = 2   // 检测通过
	ImageCheckStatusFailed = 3   // 检测失败
)

// 检测建议常量
const (
	SuggestPass = "pass"   // 通过
	SuggestReview = "review" // 需要人工审核
	SuggestRisky = "risky"   // 有风险
)
