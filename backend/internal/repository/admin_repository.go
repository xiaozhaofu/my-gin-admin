package repository

import (
	"gorm.io/gorm"

	"go_sleep_admin/internal/models"
)

type AdminRepository struct {
	db *gorm.DB
}

type AdminScope struct {
	AllData bool
	AdminID int64
	DeptIDs []int64
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) ByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Preload("Roles").Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) ByID(id int64) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Preload("Roles").First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) List(scope AdminScope) ([]models.Admin, error) {
	var admins []models.Admin
	db := r.db.Preload("Roles").Order("id asc")
	if !scope.AllData {
		switch {
		case scope.AdminID > 0:
			db = db.Where("admins.id = ?", scope.AdminID)
		case len(scope.DeptIDs) > 0:
			db = db.Where("admins.dept_id IN ?", scope.DeptIDs)
		}
	}
	err := db.Find(&admins).Error
	return admins, err
}

func (r *AdminRepository) Save(admin *models.Admin, roleIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(admin).Error; err != nil {
			return err
		}
		if roleIDs != nil {
			if err := tx.Where("admin_id = ?", admin.ID).Delete(&models.AdminRole{}).Error; err != nil {
				return err
			}
			if len(roleIDs) > 0 {
				items := make([]models.AdminRole, 0, len(roleIDs))
				for _, roleID := range roleIDs {
					items = append(items, models.AdminRole{AdminID: admin.ID, RoleID: roleID})
				}
				if err := tx.Create(&items).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
