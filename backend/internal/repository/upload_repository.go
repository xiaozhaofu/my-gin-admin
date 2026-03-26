package repository

import (
	"strings"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type UploadRepository struct {
	db *gorm.DB
}

type UploadScope struct {
	AllData bool
	AdminID int64
	DeptIDs []int64
}

func NewUploadRepository(db *gorm.DB) *UploadRepository { return &UploadRepository{db: db} }

func (r *UploadRepository) Save(item *models.UploadFile) error { return r.db.Create(item).Error }

func (r *UploadRepository) List(query dto.UploadListQuery, scope UploadScope) ([]models.UploadFile, int64, error) {
	var (
		items []models.UploadFile
		total int64
	)
	page, pageSize := query.Normalize()
	db := r.db.Model(&models.UploadFile{}).Order("id desc")
	if !scope.AllData {
		switch {
		case scope.AdminID > 0:
			db = db.Where("upload_files.admin_id = ?", scope.AdminID)
		case len(scope.DeptIDs) > 0:
			db = db.Joins("JOIN admins ON admins.id = upload_files.admin_id").Where("admins.dept_id IN ?", scope.DeptIDs)
		}
	}
	if query.OriginName != "" {
		db = db.Where("origin_name LIKE ?", "%"+query.OriginName+"%")
	}
	if query.Type != nil {
		db = db.Where("type = ?", *query.Type)
	}
	if strings.TrimSpace(query.Scene) != "" {
		db = db.Where("scene = ?", strings.TrimSpace(query.Scene))
	}
	if strings.TrimSpace(query.Provider) != "" {
		db = db.Where("provider = ?", strings.TrimSpace(query.Provider))
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *UploadRepository) Delete(id int64) error {
	return r.db.Delete(&models.UploadFile{}, id).Error
}

func (r *UploadRepository) DeleteBatch(ids []int64) error {
	return r.db.Where("id IN ?", ids).Delete(&models.UploadFile{}).Error
}
