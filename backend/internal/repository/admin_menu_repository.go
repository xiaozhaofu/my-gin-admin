package repository

import (
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type AdminMenuRepository struct {
	db *gorm.DB
}

func NewAdminMenuRepository(db *gorm.DB) *AdminMenuRepository {
	return &AdminMenuRepository{db: db}
}

func (r *AdminMenuRepository) List() ([]models.AdminMenu, error) {
	var items []models.AdminMenu
	err := r.db.Order("sort asc, id asc").Find(&items).Error
	return items, err
}

func (r *AdminMenuRepository) Save(item *models.AdminMenu) error {
	return r.db.Save(item).Error
}

func (r *AdminMenuRepository) Delete(id int64) error {
	return r.db.Delete(&models.AdminMenu{}, id).Error
}
