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

func (r *MenuRepository) Save(item *models.Menu) error { return r.db.Save(item).Error }
func (r *MenuRepository) Delete(id int64) error        { return r.db.Delete(&models.Menu{}, id).Error }
