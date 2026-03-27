package repository

import (
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository { return &MenuRepository{db: db} }

func (r *MenuRepository) List() ([]models.Menu, error) {
	var items []models.Menu
	err := r.db.Order("level asc, sort_order desc, id asc").Find(&items).Error
	return items, err
}

func (r *MenuRepository) ByID(id int64) (*models.Menu, error) {
	var item models.Menu
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MenuRepository) Save(item *models.Menu) error {
	if item.ID == 0 {
		return r.db.Create(item).Error
	}

	updates := map[string]any{
		"name":       item.Name,
		"parent_id":  item.ParentID,
		"level":      item.Level,
		"sort_order": item.SortOrder,
		"is_active":  item.IsActive,
		"page_path":  item.PagePath,
		"icon":       item.Icon,
		"admin_id":   item.AdminID,
		"updated_at": gorm.Expr("CURRENT_TIMESTAMP(3)"),
	}
	return r.db.Model(&models.Menu{}).Where("id = ?", item.ID).Updates(updates).Error
}

func (r *MenuRepository) UpdateStatus(ids []int64, isActive bool) error {
	return r.db.Model(&models.Menu{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"is_active":  isActive,
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP(3)"),
		}).Error
}

func (r *MenuRepository) Delete(id int64) error { return r.db.Delete(&models.Menu{}, id).Error }
