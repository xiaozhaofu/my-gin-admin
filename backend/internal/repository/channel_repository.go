package repository

import (
	"errors"

	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository { return &ChannelRepository{db: db} }

func (r *ChannelRepository) List() ([]models.Channel, error) {
	var items []models.Channel
	err := r.db.Order("id asc").Find(&items).Error
	return items, err
}

func (r *ChannelRepository) ByID(id int64) (*models.Channel, error) {
	var item models.Channel
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ChannelRepository) Save(item *models.Channel) error {
	if item.ID == 0 {
		return r.db.Create(item).Error
	}

	updates := map[string]any{
		"name":       item.Name,
		"code":       item.Code,
		"status":     item.Status,
		"remark":     item.Remark,
		"admin_id":   item.AdminID,
		"updated_at": gorm.Expr("CURRENT_TIMESTAMP(3)"),
	}
	return r.db.Model(&models.Channel{}).Where("id = ?", item.ID).Updates(updates).Error
}

func (r *ChannelRepository) UpdateStatus(ids []int64, status models.ChannelStatus) error {
	return r.db.Model(&models.Channel{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"status":     status,
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP(3)"),
		}).Error
}

func (r *ChannelRepository) ExistsByNameOrCode(name, code string, excludeID int64) (bool, error) {
	db := r.db.Model(&models.Channel{}).Where("name = ? OR code = ?", name, code)
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsChannelNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
