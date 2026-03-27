package service

import (
	"slices"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type MenuService struct {
	repo *repository.MenuRepository
}

func NewMenuService(repo *repository.MenuRepository) *MenuService { return &MenuService{repo: repo} }

type MenuTreeNode struct {
	models.Menu
	Children []*MenuTreeNode `json:"children,omitempty"`
}

func (s *MenuService) ListTree() ([]*MenuTreeNode, error) {
	items, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	nodeMap := map[int64]*MenuTreeNode{}
	for _, item := range items {
		it := item
		nodeMap[item.ID] = &MenuTreeNode{Menu: it}
	}
	roots := make([]*MenuTreeNode, 0)
	for _, item := range items {
		node := nodeMap[item.ID]
		if isRootMenu(item.ParentID) {
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

func (s *MenuService) Cascader() ([]*MenuTreeNode, error) {
	return s.ListTree()
}

func (s *MenuService) Save(id, adminID int64, req dto.MenuUpsertRequest) error {
	if id > 0 && req.ParentID != nil && *req.ParentID == id {
		return ErrMenuParentSelf
	}

	items, err := s.repo.List()
	if err != nil {
		return err
	}

	level := uint8(1)
	if req.ParentID != nil {
		parent, ok := findMenuByID(items, *req.ParentID)
		if !ok {
			return ErrMenuParentNotFound
		}
		if parent.Level >= 3 {
			return ErrMenuParentTooDeep
		}
		if id > 0 && isDescendantMenu(items, id, *req.ParentID) {
			return ErrMenuParentCircular
		}
		level = parent.Level + 1
	}

	item := &models.Menu{
		BaseID:    models.BaseID{ID: id},
		Name:      req.Name,
		ParentID:  req.ParentID,
		Level:     level,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		PagePath:  req.PagePath,
		Icon:      req.Icon,
		AdminID:   adminID,
	}
	if id > 0 && item.AdminID == 0 {
		existing, err := s.repo.ByID(id)
		if err != nil {
			return err
		}
		item.AdminID = existing.AdminID
	}
	return s.repo.Save(item)
}

func (s *MenuService) Delete(id int64) error {
	items, err := s.repo.List()
	if err != nil {
		return err
	}
	childExists := slices.ContainsFunc(items, func(item models.Menu) bool {
		return item.ParentID != nil && *item.ParentID == id
	})
	if childExists {
		return ErrMenuHasChildren
	}
	return s.repo.Delete(id)
}

func (s *MenuService) UpdateStatus(ids []int64, isActive bool) error {
	return s.repo.UpdateStatus(ids, isActive)
}

var ErrMenuHasChildren = &MenuServiceError{Message: "当前菜单存在子级，不能直接删除"}
var ErrMenuParentNotFound = &MenuServiceError{Message: "父级菜单不存在"}
var ErrMenuParentTooDeep = &MenuServiceError{Message: "当前仅支持三级菜单，不能继续向下添加"}
var ErrMenuParentSelf = &MenuServiceError{Message: "父级菜单不能选择自己"}
var ErrMenuParentCircular = &MenuServiceError{Message: "父级菜单不能选择自己的子级"}

type MenuServiceError struct {
	Message string
}

func (e *MenuServiceError) Error() string { return e.Message }

func isRootMenu(parentID *int64) bool {
	return parentID == nil || *parentID == 0
}

func findMenuByID(items []models.Menu, id int64) (models.Menu, bool) {
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return models.Menu{}, false
}

func isDescendantMenu(items []models.Menu, rootID, targetID int64) bool {
	childrenByParent := make(map[int64][]int64)
	for _, item := range items {
		if item.ParentID == nil {
			continue
		}
		childrenByParent[*item.ParentID] = append(childrenByParent[*item.ParentID], item.ID)
	}

	queue := append([]int64(nil), childrenByParent[rootID]...)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == targetID {
			return true
		}
		queue = append(queue, childrenByParent[current]...)
	}

	return false
}
