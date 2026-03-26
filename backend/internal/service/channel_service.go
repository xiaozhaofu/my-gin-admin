package service

import (
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type ChannelService struct {
	repo *repository.ChannelRepository
}

func NewChannelService(repo *repository.ChannelRepository) *ChannelService {
	return &ChannelService{repo: repo}
}

func (s *ChannelService) List() ([]models.Channel, error) {
	return s.repo.List()
}
