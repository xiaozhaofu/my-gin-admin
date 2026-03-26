package service

import (
	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService { return &LogService{repo: repo} }

func (s *LogService) Create(item *models.OperationLog) error { return s.repo.Create(item) }
func (s *LogService) List(query dto.OperationLogQuery) ([]models.OperationLog, int64, error) {
	return s.repo.List(query)
}
func (s *LogService) Delete(id int64) error { return s.repo.Delete(id) }
func (s *LogService) Clear() error          { return s.repo.Clear() }
