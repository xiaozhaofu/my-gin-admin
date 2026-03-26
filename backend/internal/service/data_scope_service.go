package service

import (
	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type DataScopeService struct {
	repo *repository.DataScopeRepository
}

func NewDataScopeService(repo *repository.DataScopeRepository) *DataScopeService {
	return &DataScopeService{repo: repo}
}

func (s *DataScopeService) DeptTree() ([]*dto.DeptTree, error) {
	items, err := s.repo.ListDepts()
	if err != nil {
		return nil, err
	}
	nodeMap := map[int64]*dto.DeptTree{}
	for _, item := range items {
		nodeMap[item.ID] = &dto.DeptTree{
			ID:       item.ID,
			ParentID: item.ParentID,
			Name:     item.Name,
			Code:     item.Code,
			Status:   item.Status,
		}
	}
	roots := make([]*dto.DeptTree, 0)
	for _, item := range items {
		node := nodeMap[item.ID]
		if item.ParentID == nil {
			roots = append(roots, node)
			continue
		}
		parent := nodeMap[*item.ParentID]
		if parent != nil {
			parent.Children = append(parent.Children, node)
		}
	}
	return roots, nil
}

func (s *DataScopeService) SaveDept(id int64, req dto.DeptUpsertRequest) error {
	item := &models.Dept{
		BaseID:   models.BaseID{ID: id},
		ParentID: req.ParentID,
		Name:     req.Name,
		Code:     req.Code,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Sort:     req.Sort,
		Status:   req.Status,
	}
	if item.Status == 0 {
		item.Status = 1
	}
	return s.repo.SaveDept(item)
}

func (s *DataScopeService) DeleteDept(id int64) error { return s.repo.DeleteDept(id) }

func (s *DataScopeService) Posts() ([]models.Post, error) {
	return s.repo.ListPosts()
}

func (s *DataScopeService) SavePost(id int64, req dto.PostUpsertRequest) error {
	item := &models.Post{
		BaseID: models.BaseID{ID: id},
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}
	if item.Status == 0 {
		item.Status = 1
	}
	return s.repo.SavePost(item)
}

func (s *DataScopeService) DeletePost(id int64) error { return s.repo.DeletePost(id) }
