package dao

import (
	"gorm.io/gorm"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// ImageCheckDaoImpl 图片检测数据访问实现
type ImageCheckDaoImpl struct {
	db *gorm.DB
}

// NewImageCheckDao 创建图片检测DAO实例
func NewImageCheckDao() ImageCheckDao {
	return &ImageCheckDaoImpl{db: db.GetDB()}
}

// Create 创建图片检测记录
func (d *ImageCheckDaoImpl) Create(imageCheck *model.ImageCheckModel) error {
	return d.db.Create(imageCheck).Error
}

// GetByTraceId 根据trace_id获取检测记录
func (d *ImageCheckDaoImpl) GetByTraceId(traceId string) (*model.ImageCheckModel, error) {
	var imageCheck model.ImageCheckModel
	err := d.db.Where("trace_id = ?", traceId).First(&imageCheck).Error
	if err != nil {
		return nil, err
	}
	return &imageCheck, nil
}

// UpdateStatus 更新检测状态
func (d *ImageCheckDaoImpl) UpdateStatus(traceId string, status int, suggest string, label int, prob float64, strategy string, errcode int, errmsg string) error {
	return d.db.Model(&model.ImageCheckModel{}).
		Where("trace_id = ?", traceId).
		Updates(map[string]interface{}{
			"status":   status,
			"suggest":  suggest,
			"label":    label,
			"prob":     prob,
			"strategy": strategy,
			"errcode":  errcode,
			"errmsg":   errmsg,
		}).Error
}

// GetByPostId 获取帖子的所有图片检测记录
func (d *ImageCheckDaoImpl) GetByPostId(postId int64) ([]*model.ImageCheckModel, error) {
	var imageChecks []*model.ImageCheckModel
	err := d.db.Where("post_id = ?", postId).Find(&imageChecks).Error
	return imageChecks, err
}

// GetPendingChecks 获取待检测的记录
func (d *ImageCheckDaoImpl) GetPendingChecks() ([]*model.ImageCheckModel, error) {
	var imageChecks []*model.ImageCheckModel
	err := d.db.Where("status = ?", model.ImageCheckStatusPending).Find(&imageChecks).Error
	return imageChecks, err
}

// DeleteByPostId 删除帖子的所有检测记录
func (d *ImageCheckDaoImpl) DeleteByPostId(postId int64) error {
	return d.db.Where("post_id = ?", postId).Delete(&model.ImageCheckModel{}).Error
}
