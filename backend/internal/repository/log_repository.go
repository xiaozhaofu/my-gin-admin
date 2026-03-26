package repository

import (
	"strings"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository { return &LogRepository{db: db} }

func (r *LogRepository) Create(item *models.OperationLog) error {
	return r.db.Create(item).Error
}

func (r *LogRepository) List(query dto.OperationLogQuery) ([]models.OperationLog, int64, error) {
	page, pageSize := query.Normalize()
	db := r.db.Model(&models.OperationLog{}).Order("id desc")
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Method != "" {
		db = db.Where("method = ?", strings.ToUpper(query.Method))
	}
	if query.Path != "" {
		db = db.Where("path LIKE ?", "%"+query.Path+"%")
	}
	if query.Success != nil {
		db = db.Where("success = ?", *query.Success)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.OperationLog
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *LogRepository) Delete(id int64) error {
	return r.db.Delete(&models.OperationLog{}, id).Error
}

func (r *LogRepository) Clear() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.OperationLog{}).Error
}
