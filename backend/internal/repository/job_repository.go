package repository

import (
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository { return &JobRepository{db: db} }

func (r *JobRepository) List() ([]models.SysJob, error) {
	var items []models.SysJob
	err := r.db.Order("id asc").Find(&items).Error
	return items, err
}

func (r *JobRepository) Save(item *models.SysJob) error {
	return r.db.Save(item).Error
}

func (r *JobRepository) Delete(id int64) error {
	return r.db.Delete(&models.SysJob{}, id).Error
}
