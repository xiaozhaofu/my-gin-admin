package repository

import (
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
