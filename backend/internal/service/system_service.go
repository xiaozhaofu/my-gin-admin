package service

import (
	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type SystemService struct {
	repo *repository.SystemRepository
}

func NewSystemService(repo *repository.SystemRepository) *SystemService {
	return &SystemService{repo: repo}
}

func (s *SystemService) DictTypes() ([]models.DictType, error) {
	return s.repo.ListDictTypes()
}

func (s *SystemService) SaveDictType(id int64, req dto.DictTypeUpsertRequest) error {
	item := &models.DictType{
		BaseID:   models.BaseID{ID: id},
		Name:     req.Name,
		TypeCode: req.TypeCode,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if item.Status == 0 {
		item.Status = 1
	}
	return s.repo.SaveDictType(item)
}

func (s *SystemService) SaveDictItem(id int64, req dto.DictItemUpsertRequest) error {
	item := &models.DictItem{
		BaseID:    models.BaseID{ID: id},
		TypeID:    req.TypeID,
		Label:     req.Label,
		Value:     req.Value,
		Sort:      req.Sort,
		Status:    req.Status,
		CSSClass:  req.CSSClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Remark:    req.Remark,
	}
	if item.Status == 0 {
		item.Status = 1
	}
	return s.repo.SaveDictItem(item)
}

func (s *SystemService) DeleteDictType(id int64) error { return s.repo.DeleteDictType(id) }
func (s *SystemService) DeleteDictItem(id int64) error { return s.repo.DeleteDictItem(id) }

func (s *SystemService) Configs() ([]models.SysConfig, error) {
	return s.repo.ListConfigs()
}

func (s *SystemService) SaveConfig(id int64, req dto.SysConfigUpsertRequest) error {
	item := &models.SysConfig{
		BaseID:      models.BaseID{ID: id},
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      req.Remark,
	}
	return s.repo.SaveConfig(item)
}

func (s *SystemService) DeleteConfig(id int64) error { return s.repo.DeleteConfig(id) }
