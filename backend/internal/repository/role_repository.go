package repository

import (
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository { return &RoleRepository{db: db} }

func (r *RoleRepository) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Order("id asc").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) Save(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) MenuIDs(roleID int64) ([]int64, error) {
	var links []models.RoleMenu
	if err := r.db.Where("role_id = ?", roleID).Find(&links).Error; err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(links))
	for _, link := range links {
		out = append(out, link.MenuID)
	}
	return out, nil
}
