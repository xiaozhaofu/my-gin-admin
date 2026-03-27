package service

import (
	"slices"

	"gorm.io/gorm"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/platform/auth"
	"go_sleep_admin/internal/repository"
)

type CasbinEnforcer interface {
	RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error)
	AddPolicies(rules [][]string) (bool, error)
	LoadPolicy() error
	GetPermissionsForUser(user string, domain ...string) ([][]string, error)
}

type RBACService struct {
	db            *gorm.DB
	roleRepo      *repository.RoleRepository
	adminRepo     *repository.AdminRepository
	adminMenuRepo *repository.AdminMenuRepository
	dataScopeRepo *repository.DataScopeRepository
	enforcer      CasbinEnforcer
}

func NewRBACService(db *gorm.DB, roleRepo *repository.RoleRepository, adminRepo *repository.AdminRepository, adminMenuRepo *repository.AdminMenuRepository, dataScopeRepo *repository.DataScopeRepository, enforcer CasbinEnforcer) *RBACService {
	return &RBACService{db: db, roleRepo: roleRepo, adminRepo: adminRepo, adminMenuRepo: adminMenuRepo, dataScopeRepo: dataScopeRepo, enforcer: enforcer}
}

func (s *RBACService) Roles() ([]dto.RoleDetail, error) {
	roles, err := s.roleRepo.List()
	if err != nil {
		return nil, err
	}
	out := make([]dto.RoleDetail, 0, len(roles))
	for _, role := range roles {
		menuIDs, err := s.roleRepo.MenuIDs(role.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, dto.RoleDetail{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			IsBuiltIn:   role.IsBuiltIn,
			DataScope:   role.DataScope,
			MenuIDs:     menuIDs,
		})
	}
	return out, nil
}

func (s *RBACService) SaveRole(id int64, req dto.RoleUpsertRequest) error {
	role := &models.Role{
		BaseID:      models.BaseID{ID: id},
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		DataScope:   req.DataScope,
	}
	if err := s.roleRepo.Save(role); err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", role.ID).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		if len(req.MenuIDs) > 0 {
			links := make([]models.RoleMenu, 0, len(req.MenuIDs))
			for _, menuID := range req.MenuIDs {
				links = append(links, models.RoleMenu{RoleID: role.ID, MenuID: menuID})
			}
			if err := tx.Create(&links).Error; err != nil {
				return err
			}
		}
		return syncRolePolicies(tx, s.enforcer, role)
	})
}

func (s *RBACService) Admins(claims *auth.Claims) ([]dto.AdminDetail, error) {
	scope := repository.AdminScope{}
	if claims != nil {
		admin, err := s.adminRepo.ByID(claims.AdminID)
		if err != nil {
			return nil, err
		}
		resolved, err := resolveAdminScope(admin, s.dataScopeRepo)
		if err != nil {
			return nil, err
		}
		scope = repository.AdminScope{AllData: resolved.AllData, AdminID: resolved.AdminID, DeptIDs: resolved.DeptIDs}
	}
	admins, err := s.adminRepo.List(scope)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AdminDetail, 0, len(admins))
	for _, admin := range admins {
		roleIDs := make([]int64, 0, len(admin.Roles))
		roleNames := make([]string, 0, len(admin.Roles))
		for _, role := range admin.Roles {
			roleIDs = append(roleIDs, role.ID)
			roleNames = append(roleNames, role.Name)
		}
		out = append(out, dto.AdminDetail{
			ID:       admin.ID,
			Username: admin.Username,
			Nickname: admin.Nickname,
			Phone:    admin.Phone,
			Email:    admin.Email,
			Status:   int(admin.Status),
			RoleIDs:  roleIDs,
			Roles:    roleNames,
			DeptID:   admin.DeptID,
			PostID:   admin.PostID,
		})
	}
	return out, nil
}

func (s *RBACService) SaveAdmin(id int64, passwordHash string, req dto.AdminUpsertRequest) error {
	status := models.AdminStatus(req.Status)
	if status == 0 {
		status = models.AdminStatusNormal
	}
	admin := &models.Admin{
		BaseID:   models.BaseID{ID: id},
		Username: req.Username,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: passwordHash,
		Status:   status,
		DeptID:   req.DeptID,
		PostID:   req.PostID,
	}
	if id > 0 && passwordHash == "" {
		existing, err := s.adminRepo.ByID(id)
		if err != nil {
			return err
		}
		admin.Password = existing.Password
	}
	return s.adminRepo.Save(admin, req.RoleIDs)
}

func syncRolePolicies(db *gorm.DB, enforcer CasbinEnforcer, role *models.Role) error {
	if _, err := enforcer.RemoveFilteredPolicy(0, role.Code); err != nil {
		return err
	}

	var menus []models.AdminMenu
	err := db.Table("admin_menus").
		Select("admin_menus.*").
		Joins("join role_menus on role_menus.menu_id = admin_menus.id").
		Where("role_menus.role_id = ?", role.ID).
		Where("admin_menus.api_path != '' and admin_menus.method != ''").
		Find(&menus).Error
	if err != nil {
		return err
	}

	policies := make([][]string, 0, len(menus))
	seen := make(map[string]struct{}, len(menus))
	for _, menu := range menus {
		for _, rule := range expandAdminMenuPolicies(role.Code, menu) {
			key := rule[0] + "\x00" + rule[1] + "\x00" + rule[2]
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			policies = append(policies, rule)
		}
	}
	if len(policies) > 0 {
		if _, err := enforcer.AddPolicies(policies); err != nil {
			return err
		}
	}
	return enforcer.LoadPolicy()
}

func expandAdminMenuPolicies(roleCode string, menu models.AdminMenu) [][]string {
	if menu.APIPath == "" || menu.Method == "" {
		return nil
	}

	rules := [][]string{{roleCode, menu.APIPath, menu.Method}}

	// 内容菜单页的树接口和级联接口属于同一组读取权限。
	// 历史上后台菜单只登记了 /menus/tree，这里额外补出 /menus/cascader，
	// 避免文章创建页因缺少读取级联菜单权限而拿不到数据。
	if menu.APIPath == "/api/v1/menus/tree" && menu.Method == "GET" {
		rules = append(rules, []string{roleCode, "/api/v1/menus/cascader", "GET"})
	}
	if menu.APIPath == "/api/v1/articles" && menu.Method == "GET" {
		rules = append(rules, []string{roleCode, "/api/v1/channels", "GET"})
	}
	if menu.APIPath == "/api/v1/articles" && menu.Method == "POST" {
		rules = append(rules, []string{roleCode, "/api/v1/channels", "GET"})
		rules = append(rules, []string{roleCode, "/api/v1/menus/cascader", "GET"})
	}
	if menu.APIPath == "/api/v1/articles/batch" && menu.Method == "POST" {
		rules = append(rules, []string{roleCode, "/api/v1/channels", "GET"})
		rules = append(rules, []string{roleCode, "/api/v1/menus/cascader", "GET"})
	}
	if menu.APIPath == "/api/v1/orders" && menu.Method == "GET" {
		rules = append(rules, []string{roleCode, "/api/v1/orders/:id", "GET"})
		rules = append(rules, []string{roleCode, "/api/v1/orders/export", "GET"})
	}
	if menu.APIPath == "/api/v1/uploads" && menu.Method == "GET" {
		rules = append(rules, []string{roleCode, "/api/v1/uploads", "DELETE"})
		rules = append(rules, []string{roleCode, "/api/v1/uploads/:id", "DELETE"})
	}

	return rules
}

func (s *RBACService) AdminMenuTree() ([]*dto.AdminMenuTree, error) {
	items, err := s.adminMenuRepo.List()
	if err != nil {
		return nil, err
	}
	nodeMap := map[int64]*dto.AdminMenuTree{}
	for _, item := range items {
		nodeMap[item.ID] = &dto.AdminMenuTree{
			ID:         item.ID,
			ParentID:   item.ParentID,
			Title:      item.Title,
			Name:       item.Name,
			Path:       item.Path,
			Component:  item.Component,
			Icon:       item.Icon,
			Permission: item.Permission,
			Type:       item.Type,
			Sort:       item.Sort,
			Hidden:     item.Hidden,
			KeepAlive:  item.KeepAlive,
			Method:     item.Method,
			APIPath:    item.APIPath,
		}
	}
	roots := make([]*dto.AdminMenuTree, 0)
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
	for _, root := range roots {
		sortAdminMenuTree(root)
	}
	return roots, nil
}

func (s *RBACService) SaveAdminMenu(id int64, req dto.AdminMenuUpsertRequest) error {
	item := &models.AdminMenu{
		BaseID:     models.BaseID{ID: id},
		ParentID:   req.ParentID,
		Title:      req.Title,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Permission: req.Permission,
		Type:       req.Type,
		Sort:       req.Sort,
		Hidden:     req.Hidden,
		KeepAlive:  req.KeepAlive,
		Method:     req.Method,
		APIPath:    req.APIPath,
	}
	if item.Type == 0 {
		item.Type = 1
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.adminMenuRepo.Save(item); err != nil {
			return err
		}
		return s.resyncPolicies()
	})
}

func (s *RBACService) DeleteAdminMenu(id int64) error {
	items, err := s.adminMenuRepo.List()
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ParentID != nil && *item.ParentID == id {
			return &MenuServiceError{Message: "当前后台菜单存在子级，不能直接删除"}
		}
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", id).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := s.adminMenuRepo.Delete(id); err != nil {
			return err
		}
		return s.resyncPolicies()
	})
}

func (s *RBACService) resyncPolicies() error {
	roles, err := s.roleRepo.List()
	if err != nil {
		return err
	}
	for _, role := range roles {
		if err := syncRolePolicies(s.db, s.enforcer, &role); err != nil {
			return err
		}
	}
	return nil
}

func sortAdminMenuTree(node *dto.AdminMenuTree) {
	slices.SortFunc(node.Children, func(a, b *dto.AdminMenuTree) int {
		switch {
		case a.ID < b.ID:
			return -1
		case a.ID > b.ID:
			return 1
		default:
			return 0
		}
	})
	for _, child := range node.Children {
		sortAdminMenuTree(child)
	}
}
