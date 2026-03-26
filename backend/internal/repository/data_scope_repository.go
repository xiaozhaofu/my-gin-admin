package repository

import (
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type DataScopeRepository struct {
	db *gorm.DB
}

func NewDataScopeRepository(db *gorm.DB) *DataScopeRepository { return &DataScopeRepository{db: db} }

func (r *DataScopeRepository) ListDepts() ([]models.Dept, error) {
	var items []models.Dept
	err := r.db.Order("sort asc, id asc").Find(&items).Error
	return items, err
}

func (r *DataScopeRepository) SaveDept(item *models.Dept) error {
	return r.db.Save(item).Error
}

func (r *DataScopeRepository) DeleteDept(id int64) error {
	return r.db.Delete(&models.Dept{}, id).Error
}

func (r *DataScopeRepository) ListPosts() ([]models.Post, error) {
	var items []models.Post
	err := r.db.Order("sort asc, id asc").Find(&items).Error
	return items, err
}

func (r *DataScopeRepository) SavePost(item *models.Post) error {
	return r.db.Save(item).Error
}

func (r *DataScopeRepository) DeletePost(id int64) error {
	return r.db.Delete(&models.Post{}, id).Error
}
