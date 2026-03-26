package service

import (
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type AdminScope struct {
	AllData bool
	AdminID int64
	DeptIDs []int64
}

func resolveAdminScope(admin *models.Admin, dataScopeRepo *repository.DataScopeRepository) (AdminScope, error) {
	scopeLevel := int8(4)
	for _, role := range admin.Roles {
		switch role.DataScope {
		case 1:
			return AdminScope{AllData: true}, nil
		case 3:
			if scopeLevel > 3 {
				scopeLevel = 3
			}
		case 2:
			if scopeLevel > 2 {
				scopeLevel = 2
			}
		case 4:
			if scopeLevel > 4 {
				scopeLevel = 4
			}
		}
	}

	if admin.DeptID == nil {
		return AdminScope{AdminID: admin.ID}, nil
	}

	switch scopeLevel {
	case 2:
		return AdminScope{DeptIDs: []int64{*admin.DeptID}}, nil
	case 3:
		depts, err := dataScopeRepo.ListDepts()
		if err != nil {
			return AdminScope{}, err
		}
		return AdminScope{DeptIDs: collectDescendantDeptIDs(*admin.DeptID, depts)}, nil
	default:
		return AdminScope{AdminID: admin.ID}, nil
	}
}

func collectDescendantDeptIDs(rootID int64, depts []models.Dept) []int64 {
	ids := []int64{rootID}
	changed := true
	for changed {
		changed = false
		for _, dept := range depts {
			if dept.ParentID == nil {
				continue
			}
			if containsInt64(ids, *dept.ParentID) && !containsInt64(ids, dept.ID) {
				ids = append(ids, dept.ID)
				changed = true
			}
		}
	}
	return ids
}

func containsInt64(items []int64, target int64) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
