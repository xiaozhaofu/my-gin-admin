package service

import (
	"strings"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(repo *repository.JobRepository) *JobService { return &JobService{repo: repo} }

func (s *JobService) List() ([]models.SysJob, error) {
	return s.repo.List()
}

func (s *JobService) Save(id int64, req dto.SysJobUpsertRequest) error {
	item := &models.SysJob{
		BaseID:         models.BaseID{ID: id},
		Name:           req.Name,
		Executor:       strings.TrimSpace(req.Executor),
		JobKey:         req.JobKey,
		CronExpression: req.CronExpr,
		CronExpr:       req.CronExpr,
		Target:         req.Target,
		Status:         req.Status,
		Concurrent:     req.Concurrent,
		Remark:         req.Remark,
	}
	if item.Executor == "" {
		item.Executor = "system"
	}
	if item.Status == 0 {
		item.Status = 1
	}
	return s.repo.Save(item)
}

func (s *JobService) Delete(id int64) error { return s.repo.Delete(id) }
