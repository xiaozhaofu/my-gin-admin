package repository

import (
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) *SystemRepository { return &SystemRepository{db: db} }

func (r *SystemRepository) ListDictTypes() ([]models.DictType, error) {
	var items []models.DictType
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort asc, id asc")
	}).Order("id asc").Find(&items).Error
	return items, err
}

func (r *SystemRepository) SaveDictType(item *models.DictType) error {
	return r.db.Save(item).Error
}

func (r *SystemRepository) SaveDictItem(item *models.DictItem) error {
	return r.db.Save(item).Error
}

func (r *SystemRepository) DeleteDictType(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("type_id = ?", id).Delete(&models.DictItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.DictType{}, id).Error
	})
}

func (r *SystemRepository) DeleteDictItem(id int64) error {
	return r.db.Delete(&models.DictItem{}, id).Error
}

func (r *SystemRepository) ListConfigs() ([]models.SysConfig, error) {
	var items []models.SysConfig
	err := r.db.Order("id asc").Find(&items).Error
	return items, err
}

func (r *SystemRepository) SaveConfig(item *models.SysConfig) error {
	return r.db.Save(item).Error
}

func (r *SystemRepository) DeleteConfig(id int64) error {
	return r.db.Delete(&models.SysConfig{}, id).Error
}
