package bootstrap

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/platform/auth"
	v2config "go_sleep_admin/internal/platform/config"
)

func seed(cfg *v2config.Config, db *gorm.DB, enforcer *casbin.Enforcer, _ *auth.JWTManager) error {
	resetBuiltinData := shouldResetBuiltinData(cfg)

	passwordAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	password123456, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	roles := []models.Role{
		{Name: "超级管理员", Code: "admin", Description: "拥有全部权限", IsBuiltIn: true},
		{Name: "内容运营", Code: "content_editor", Description: "管理文章和资源", IsBuiltIn: true},
		{Name: "审核员", Code: "reviewer", Description: "审核内容", IsBuiltIn: true},
	}
	for _, role := range roles {
		tx := db.Where("code = ?", role.Code)
		if resetBuiltinData {
			tx = tx.Assign(map[string]any{
				"name":        role.Name,
				"description": role.Description,
				"is_built_in": role.IsBuiltIn,
			})
		}
		if err := tx.FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	admins := []models.Admin{
		{Username: "admin", Nickname: "超级管理员", Password: string(passwordAdmin), Status: models.AdminStatusNormal},
		{Username: "operator", Nickname: "内容运营", Password: string(password123456), Status: models.AdminStatusNormal},
		{Username: "reviewer", Nickname: "审核员", Password: string(password123456), Status: models.AdminStatusNormal},
	}
	for _, item := range admins {
		tx := db.Where("username = ?", item.Username)
		if resetBuiltinData {
			tx = tx.Assign(map[string]any{
				"nickname": item.Nickname,
				"password": item.Password,
				"status":   item.Status,
			})
		}
		if err := tx.FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}

	var adminRole, editorRole, reviewerRole models.Role
	if err := db.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("code = ?", "content_editor").First(&editorRole).Error; err != nil {
		return err
	}
	if err := db.Where("code = ?", "reviewer").First(&reviewerRole).Error; err != nil {
		return err
	}
	var rootAdmin, editorAdmin, reviewerAdmin models.Admin
	if err := db.Where("username = ?", "admin").First(&rootAdmin).Error; err != nil {
		return err
	}
	if err := db.Where("username = ?", "operator").First(&editorAdmin).Error; err != nil {
		return err
	}
	if err := db.Where("username = ?", "reviewer").First(&reviewerAdmin).Error; err != nil {
		return err
	}

	roleLinks := []models.AdminRole{
		{AdminID: rootAdmin.ID, RoleID: adminRole.ID},
		{AdminID: editorAdmin.ID, RoleID: editorRole.ID},
		{AdminID: reviewerAdmin.ID, RoleID: reviewerRole.ID},
	}
	for _, link := range roleLinks {
		if err := db.Where("admin_id = ? and role_id = ?", link.AdminID, link.RoleID).FirstOrCreate(&link).Error; err != nil {
			return err
		}
	}

	rootDept := models.Dept{Name: "总部", Code: "hq", Leader: "admin", Status: 1}
	if err := db.Where("code = ?", rootDept.Code).FirstOrCreate(&rootDept).Error; err != nil {
		return err
	}
	contentDept := models.Dept{Name: "内容中心", Code: "content-center", ParentID: &rootDept.ID, Leader: "operator", Status: 1}
	if err := db.Where("code = ?", contentDept.Code).FirstOrCreate(&contentDept).Error; err != nil {
		return err
	}
	postAdmin := models.Post{Name: "系统管理员", Code: "sys_admin", Sort: 1, Status: 1}
	if err := db.Where("code = ?", postAdmin.Code).FirstOrCreate(&postAdmin).Error; err != nil {
		return err
	}
	postEditor := models.Post{Name: "内容运营", Code: "content_editor", Sort: 2, Status: 1}
	if err := db.Where("code = ?", postEditor.Code).FirstOrCreate(&postEditor).Error; err != nil {
		return err
	}
	_ = db.Model(&models.Admin{}).Where("id = ?", rootAdmin.ID).Updates(map[string]any{"dept_id": rootDept.ID, "post_id": postAdmin.ID}).Error
	_ = db.Model(&models.Admin{}).Where("id = ?", editorAdmin.ID).Updates(map[string]any{"dept_id": contentDept.ID, "post_id": postEditor.ID}).Error
	_ = db.Model(&models.Admin{}).Where("id = ?", reviewerAdmin.ID).Updates(map[string]any{"dept_id": contentDept.ID, "post_id": postEditor.ID}).Error

	menus := []models.AdminMenu{
		{Title: "工作台", Name: "dashboard", Path: "/home", Component: "home/home", Icon: "icon-home", Type: 1, Sort: 1},
		{Title: "文章管理", Name: "article-list", Path: "/articles", Component: "biz/articles/index", Icon: "icon-file", Permission: "article:list", Type: 1, Sort: 10, APIPath: "/api/v1/articles", Method: "GET"},
		{Title: "批量新增", Name: "article-batch-create", Path: "/articles/batch", Component: "biz/articles/batch", Icon: "icon-folder-add", Permission: "article:batch_create", Type: 1, Sort: 11, APIPath: "/api/v1/articles/batch", Method: "POST"},
		{Title: "文章新增", Name: "article-create", Path: "/articles/create", Component: "biz/articles/form", Hidden: true, Permission: "article:create", Type: 2, Sort: 12, APIPath: "/api/v1/articles", Method: "POST"},
		{Title: "文章编辑", Name: "article-update", Path: "/articles/:id", Component: "biz/articles/form", Hidden: true, Permission: "article:update", Type: 2, Sort: 13, APIPath: "/api/v1/articles/:id", Method: "PUT"},
		{Title: "内容菜单", Name: "content-menu", Path: "/content-menus", Component: "biz/menus/index", Icon: "icon-apps", Permission: "menu:list", Type: 1, Sort: 20, APIPath: "/api/v1/menus/tree", Method: "GET"},
		{Title: "渠道管理", Name: "channel-list", Path: "/channels", Component: "biz/channels/index", Icon: "icon-list", Permission: "channel:list", Type: 1, Sort: 25, APIPath: "/api/v1/channels", Method: "GET"},
		{Title: "渠道新增", Name: "channel-create", Path: "/channels/create", Component: "biz/channels/index", Hidden: true, Permission: "channel:create", Type: 2, Sort: 26, APIPath: "/api/v1/channels", Method: "POST"},
		{Title: "渠道编辑", Name: "channel-update", Path: "/channels/:id", Component: "biz/channels/index", Hidden: true, Permission: "channel:update", Type: 2, Sort: 27, APIPath: "/api/v1/channels/:id", Method: "PUT"},
		{Title: "渠道状态", Name: "channel-status", Path: "/channels/status", Component: "biz/channels/index", Hidden: true, Permission: "channel:status", Type: 2, Sort: 28, APIPath: "/api/v1/channels/status", Method: "PUT"},
		{Title: "资源上传", Name: "upload-file", Path: "/uploads", Component: "biz/uploads/index", Icon: "icon-upload", Permission: "upload:list", Type: 1, Sort: 30, APIPath: "/api/v1/uploads", Method: "GET"},
		{Title: "系统管理", Name: "system", Path: "/system", Component: "Layout", Icon: "icon-settings", Type: 1, Sort: 90},
		{Title: "后台菜单", Name: "admin-menu-list", Path: "/system/admin-menus", Component: "biz/system/admin-menus", Permission: "admin_menu:list", Type: 1, Sort: 90, APIPath: "/api/v1/admin-menus/tree", Method: "GET"},
		{Title: "字典管理", Name: "dict-list", Path: "/system/dicts", Component: "biz/system/dicts", Permission: "dict:list", Type: 1, Sort: 93, APIPath: "/api/v1/dict-types", Method: "GET"},
		{Title: "系统参数", Name: "config-list", Path: "/system/configs", Component: "biz/system/configs", Permission: "config:list", Type: 1, Sort: 94, APIPath: "/api/v1/sys-configs", Method: "GET"},
		{Title: "操作日志", Name: "operation-log-list", Path: "/system/operation-logs", Component: "biz/system/operation-logs", Permission: "operation_log:list", Type: 1, Sort: 95, APIPath: "/api/v1/operation-logs", Method: "GET"},
		{Title: "登录日志", Name: "login-log-list", Path: "/system/login-logs", Component: "biz/system/login-logs", Permission: "login_log:list", Type: 1, Sort: 96, APIPath: "/api/v1/login-logs", Method: "GET"},
		{Title: "在线用户", Name: "online-user-list", Path: "/system/online-users", Component: "biz/system/online-users", Permission: "online_user:list", Type: 1, Sort: 97, APIPath: "/api/v1/online-sessions", Method: "GET"},
		{Title: "API 管理", Name: "api-list", Path: "/system/apis", Component: "biz/system/apis", Permission: "api:list", Type: 1, Sort: 98, APIPath: "/api/v1/admin-menus/tree", Method: "GET"},
		{Title: "定时任务", Name: "job-list", Path: "/system/jobs", Component: "biz/system/jobs", Permission: "job:list", Type: 1, Sort: 99, APIPath: "/api/v1/jobs", Method: "GET"},
		{Title: "部门管理", Name: "dept-list", Path: "/system/depts", Component: "biz/system/depts", Permission: "dept:list", Type: 1, Sort: 100, APIPath: "/api/v1/depts/tree", Method: "GET"},
		{Title: "岗位管理", Name: "post-list", Path: "/system/posts", Component: "biz/system/posts", Permission: "post:list", Type: 1, Sort: 101, APIPath: "/api/v1/posts", Method: "GET"},
		{Title: "管理员", Name: "admin-list", Path: "/system/admins", Component: "biz/system/admins", Permission: "admin:list", Type: 1, Sort: 91, APIPath: "/api/v1/admins", Method: "GET"},
		{Title: "角色权限", Name: "role-list", Path: "/system/roles", Component: "biz/system/roles", Permission: "role:list", Type: 1, Sort: 92, APIPath: "/api/v1/roles", Method: "GET"},
	}

	for idx := range menus {
		menu := menus[idx]
		tx := db.Where("name = ?", menu.Name)
		if resetBuiltinData {
			tx = tx.Assign(menu)
		}
		err := tx.FirstOrCreate(&menu).Error
		if err != nil {
			return err
		}
		menus[idx] = menu
	}

	var editorMenuIDs []int64
	var reviewerMenuIDs []int64
	for _, menu := range menus {
		switch {
		case strings.HasPrefix(menu.Name, "article") || strings.Contains(menu.Name, "upload") || strings.Contains(menu.Name, "content-menu"):
			editorMenuIDs = append(editorMenuIDs, menu.ID)
		case strings.Contains(menu.Name, "article-list") || strings.Contains(menu.Name, "content-menu"):
			reviewerMenuIDs = append(reviewerMenuIDs, menu.ID)
		}
	}
	if resetBuiltinData {
		if err := replaceRoleMenus(db, adminRole.ID, collectMenuIDs(menus)); err != nil {
			return err
		}
		if err := replaceRoleMenus(db, editorRole.ID, editorMenuIDs); err != nil {
			return err
		}
		if err := replaceRoleMenus(db, reviewerRole.ID, reviewerMenuIDs); err != nil {
			return err
		}
	}

	if err := syncPolicies(db, enforcer); err != nil {
		return err
	}

	channel := models.Channel{Name: "默认渠道", Code: "default", Status: models.ChannelStatusNormal, Remark: "初始化渠道", AdminID: rootAdmin.ID}
	if err := db.Where("name = ?", channel.Name).FirstOrCreate(&channel).Error; err != nil {
		return err
	}

	dictType := models.DictType{Name: "文章状态", TypeCode: "article_status", Status: 1, Remark: "文章状态字典"}
	if err := db.Where("type_code = ?", dictType.TypeCode).FirstOrCreate(&dictType).Error; err != nil {
		return err
	}
	dictItems := []models.DictItem{
		{TypeID: dictType.ID, Label: "正常", Value: "1", Sort: 1, Status: 1, ListClass: "success"},
		{TypeID: dictType.ID, Label: "隐藏", Value: "2", Sort: 2, Status: 1, ListClass: "warning"},
		{TypeID: dictType.ID, Label: "待审核", Value: "3", Sort: 3, Status: 1, ListClass: "processing"},
		{TypeID: dictType.ID, Label: "已驳回", Value: "4", Sort: 4, Status: 1, ListClass: "danger"},
	}
	for _, item := range dictItems {
		if err := db.Where("type_id = ? and value = ?", item.TypeID, item.Value).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}

	sysConfigs := []models.SysConfig{
		{ConfigName: "站点名称", ConfigKey: "site.name", ConfigValue: "Sleep Admin", ConfigType: 0, Remark: "前台展示名称"},
		{ConfigName: "上传驱动", ConfigKey: "upload.driver", ConfigValue: "local", ConfigType: 0, Remark: "local 或 oss"},
	}
	for _, item := range sysConfigs {
		if err := db.Where("config_key = ?", item.ConfigKey).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}

	job := models.SysJob{Name: "刷新文章缓存", Executor: "system", JobKey: "article.cache.refresh", CronExpression: "0 */10 * * * *", CronExpr: "0 */10 * * * *", Target: "article:refresh", Status: 1, Concurrent: false, Remark: "演示任务，仅配置不执行"}
	if err := db.Where("job_key = ?", job.JobKey).FirstOrCreate(&job).Error; err != nil {
		return err
	}
	return nil
}

func shouldResetBuiltinData(cfg *v2config.Config) bool {
	if cfg == nil {
		return true
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Env)) {
	case "prod", "pro", "production":
		return false
	default:
		return true
	}
}

func replaceRoleMenus(db *gorm.DB, roleID int64, menuIDs []int64) error {
	if err := db.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
		return err
	}
	for _, menuID := range menuIDs {
		link := models.RoleMenu{RoleID: roleID, MenuID: menuID}
		if err := db.Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}

func collectMenuIDs(items []models.AdminMenu) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID)
	}
	return out
}

func syncPolicies(db *gorm.DB, enforcer *casbin.Enforcer) error {
	enforcer.ClearPolicy()
	var roles []models.Role
	if err := db.Find(&roles).Error; err != nil {
		return err
	}
	for _, role := range roles {
		if err := syncRolePoliciesForSeed(db, enforcer, &role); err != nil {
			return err
		}
	}
	return enforcer.LoadPolicy()
}

func syncRolePoliciesForSeed(db *gorm.DB, enforcer *casbin.Enforcer, role *models.Role) error {
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

	if len(menus) == 0 {
		return nil
	}

	policies := make([][]string, 0, len(menus))
	seen := make(map[string]struct{}, len(menus))
	for _, menu := range menus {
		for _, rule := range expandSeedAdminMenuPolicies(role.Code, menu) {
			key := rule[0] + "\x00" + rule[1] + "\x00" + rule[2]
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			policies = append(policies, rule)
		}
	}
	if _, err := enforcer.AddPolicies(policies); err != nil {
		return err
	}
	return nil
}

func expandSeedAdminMenuPolicies(roleCode string, menu models.AdminMenu) [][]string {
	if menu.APIPath == "" || menu.Method == "" {
		return nil
	}

	rules := [][]string{{roleCode, menu.APIPath, menu.Method}}
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
	if menu.APIPath == "/api/v1/uploads" && menu.Method == "GET" {
		rules = append(rules, []string{roleCode, "/api/v1/uploads", "DELETE"})
		rules = append(rules, []string{roleCode, "/api/v1/uploads/:id", "DELETE"})
	}

	return rules
}

func buildCasbinModel() string {
	return `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`
}
