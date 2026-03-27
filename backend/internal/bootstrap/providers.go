package bootstrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	casbinadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gtkit/json"
	"github.com/gtkit/logger"
	"github.com/gtkit/verify"
	"gorm.io/gorm"

	"go_sleep_admin/internal/middleware"
	authmodule "go_sleep_admin/internal/module/auth/transport/http"
	commonmodule "go_sleep_admin/internal/module/common/transport/http"
	contentmodule "go_sleep_admin/internal/module/content/transport/http"
	rbacmodule "go_sleep_admin/internal/module/rbac/transport/http"
	systemmodule "go_sleep_admin/internal/module/system/transport/http"
	legacyfile "go_sleep_admin/internal/pkg/file"
	legacyjwtauth "go_sleep_admin/internal/pkg/jwtauth"
	legacylog "go_sleep_admin/internal/pkg/log"
	"go_sleep_admin/internal/pkg/news"
	platformauth "go_sleep_admin/internal/platform/auth"
	v2config "go_sleep_admin/internal/platform/config"
	legacydao "go_sleep_admin/internal/platform/data"
	"go_sleep_admin/internal/platform/storage"
	v2router "go_sleep_admin/internal/router"
	"go_sleep_admin/internal/runtime/resource"
)

var (
	initLogWithConfig  = legacylog.InitWithConfig
	initNewsWithConfig = news.InitWithConfig
	initFileWithConfig = legacyfile.InitWithConfig
	initDAOWithConfig  = legacydao.NewWithConfig
	initJWTWithConfig  = legacyjwtauth.NewWithConfig
	currentJWTIssuer   = legacyjwtauth.JwtEd25519
	initVerify         = verify.New
	checkJSONRuntime   = json.CheckJSON
)

type providerSet struct {
	httpConfig HTTPConfig
	httpServer httpServer
	workers    *WorkerManager
}

func newProviders(runtime *Runtime) (*providerSet, error) {
	db := runtime.DB()
	if db == nil || db.Mdb() == nil {
		return nil, fmt.Errorf("legacy dao is not initialized")
	}

	cfg := runtime.Config()
	mysqlDB := db.Mdb()

	if err := autoMigrate(mysqlDB); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	if err := ensureArticleCoverSchema(mysqlDB); err != nil {
		return nil, fmt.Errorf("ensure article cover schema: %w", err)
	}
	if err := backfillArticleCover(mysqlDB); err != nil {
		return nil, fmt.Errorf("backfill article cover: %w", err)
	}
	if err := ensureArticleMenuSchema(mysqlDB); err != nil {
		return nil, fmt.Errorf("ensure article_menus schema: %w", err)
	}
	if err := backfillArticleMenus(mysqlDB); err != nil {
		return nil, fmt.Errorf("backfill article menus: %w", err)
	}
	if err := ensureUploadFileSchema(mysqlDB); err != nil {
		return nil, fmt.Errorf("ensure upload_files schema: %w", err)
	}
	if err := backfillUploadProvider(mysqlDB); err != nil {
		return nil, fmt.Errorf("backfill upload provider: %w", err)
	}

	enforcer, err := newEnforcer(mysqlDB)
	if err != nil {
		return nil, fmt.Errorf("init casbin enforcer: %w", err)
	}

	jwtManager := platformauth.NewJWTManager(
		runtime.TokenIssuer(),
		time.Duration(cfg.JWT.Timeout)*time.Second,
		time.Duration(cfg.JWT.RefreshTimeout)*time.Second,
	)

	uploader, err := storage.NewUploadGatewayWithConfig(cfg.Upload)
	if err != nil {
		return nil, fmt.Errorf("init uploader: %w", err)
	}

	if err := seed(cfg, mysqlDB, enforcer, jwtManager); err != nil {
		return nil, fmt.Errorf("seed database: %w", err)
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load casbin policy: %w", err)
	}

	components := newHTTPComponents(mysqlDB, enforcer, jwtManager, uploader)
	authMiddleware := middleware.JWT(jwtManager, components.Services.Session)
	permissionMiddleware := middleware.Casbin(enforcer)
	operationMiddleware := middleware.OperationLogger(components.Services.Log)

	modules := []v2router.Module{
		authmodule.NewModule(components.Handlers.Auth, authMiddleware, operationMiddleware),
		commonmodule.NewModule(components.Handlers.Module, components.Handlers.Dashboard, authMiddleware, operationMiddleware),
		rbacmodule.NewModule(components.Handlers.RBAC, authMiddleware, permissionMiddleware, operationMiddleware),
		systemmodule.NewModule(
			components.Handlers.System,
			components.Handlers.Log,
			components.Handlers.Session,
			components.Handlers.Job,
			components.Handlers.DataScope,
			authMiddleware,
			permissionMiddleware,
			operationMiddleware,
		),
		contentmodule.NewModule(
			components.Handlers.Menu,
			components.Handlers.Article,
			components.Handlers.Channel,
			components.Handlers.Order,
			components.Handlers.Upload,
			authMiddleware,
			permissionMiddleware,
			operationMiddleware,
		),
	}

	engine := v2router.InitRouterWithModules(cfg.Env, modules)
	engine.Static(cfg.Upload.PublicPath, resource.ResolvePath(cfg.Upload.LocalDir))

	httpConfig := newHTTPConfig(cfg)
	return &providerSet{
		httpConfig: httpConfig,
		httpServer: newHTTPServer(engine, httpConfig),
		workers:    NewWorkerManager(),
	}, nil
}

func initLegacyRuntime(cfg *v2config.Config) error {
	initLogWithConfig(cfg)
	initNewsWithConfig(cfg)
	initFileWithConfig(cfg)
	initDAOWithConfig(cfg)
	if _, err := initJWTWithConfig(cfg); err != nil {
		return fmt.Errorf("init jwt: %w", err)
	}
	initVerify()
	checkJSONRuntime()

	logger.Infof("runtime initialized with env=%s", cfg.Env)
	return nil
}

func closeLegacyRuntime() {
	legacydao.DBClose()
	logger.Sync()
}

func newEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := casbinadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	model, err := casbinmodel.NewModelFromString(buildCasbinModel())
	if err != nil {
		return nil, err
	}

	return casbin.NewEnforcer(model, adapter)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
	//&models.Admin{},
	//&models.Role{},
	//&models.AdminRole{},
	//&models.AdminMenu{},
	//&models.RoleMenu{},
	//&models.Dept{},
	//&models.Post{},
	//&models.DictType{},
	//&models.DictItem{},
	//&models.SysConfig{},
	//&models.OperationLog{},
	//&models.LoginLog{},
	//&models.OnlineSession{},
	//&models.SysJob{},
	)
}

func newNoopWorker(_ context.Context) error { return nil }

func ensureArticleCoverSchema(db *gorm.DB) error {
	type columnSpec struct {
		name       string
		definition string
	}

	specs := []columnSpec{
		{name: "cover_large", definition: "ALTER TABLE articles ADD COLUMN cover_large varchar(255) NOT NULL DEFAULT ''"},
		{name: "cover_medium", definition: "ALTER TABLE articles ADD COLUMN cover_medium varchar(255) NOT NULL DEFAULT ''"},
		{name: "cover_small", definition: "ALTER TABLE articles ADD COLUMN cover_small varchar(255) NOT NULL DEFAULT ''"},
	}

	for _, spec := range specs {
		exists, err := hasColumn(db, "articles", spec.name)
		if err != nil {
			return err
		}
		if !exists {
			if err := db.Exec(spec.definition).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func backfillArticleCover(db *gorm.DB) error {
	hasLegacyCover, err := hasColumn(db, "articles", "cover")
	if err != nil {
		return err
	}
	if !hasLegacyCover {
		return nil
	}

	updateSQL := `
UPDATE articles
SET cover_large = CASE WHEN cover_large = '' OR cover_large IS NULL THEN cover ELSE cover_large END,
    cover_medium = CASE WHEN cover_medium = '' OR cover_medium IS NULL THEN cover ELSE cover_medium END,
    cover_small = CASE WHEN cover_small = '' OR cover_small IS NULL THEN cover ELSE cover_small END
WHERE cover IS NOT NULL AND cover <> ''
`

	return db.Exec(updateSQL).Error
}

func ensureArticleMenuSchema(db *gorm.DB) error {
	createTableSQL := `
CREATE TABLE IF NOT EXISTS article_menus (
  article_id BIGINT UNSIGNED NOT NULL,
  menu_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (article_id, menu_id),
  KEY idx_article_menus_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
`

	return db.Exec(createTableSQL).Error
}

func backfillArticleMenus(db *gorm.DB) error {
	backfillSQL := `
INSERT INTO article_menus (article_id, menu_id)
SELECT a.id, a.menu_id
FROM articles a
WHERE a.menu_id > 0
  AND NOT EXISTS (
    SELECT 1
    FROM article_menus am
    WHERE am.article_id = a.id
      AND am.menu_id = a.menu_id
  )
`

	return db.Exec(backfillSQL).Error
}

func ensureUploadFileSchema(db *gorm.DB) error {
	type columnSpec struct {
		name       string
		definition string
		indexSQL   string
	}

	specs := []columnSpec{
		{name: "scene", definition: "ALTER TABLE upload_files ADD COLUMN scene varchar(100) NOT NULL DEFAULT 'misc'", indexSQL: "CREATE INDEX idx_upload_files_scene ON upload_files(scene)"},
		{name: "provider", definition: "ALTER TABLE upload_files ADD COLUMN provider varchar(32) NOT NULL DEFAULT 'local'", indexSQL: "CREATE INDEX idx_upload_files_provider ON upload_files(provider)"},
		{name: "created_at", definition: "ALTER TABLE upload_files ADD COLUMN created_at datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)"},
		{name: "updated_at", definition: "ALTER TABLE upload_files ADD COLUMN updated_at datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)"},
		{name: "deleted_at", definition: "ALTER TABLE upload_files ADD COLUMN deleted_at datetime(3) NULL", indexSQL: "CREATE INDEX idx_upload_files_deleted_at ON upload_files(deleted_at)"},
	}

	for _, spec := range specs {
		exists, err := hasColumn(db, "upload_files", spec.name)
		if err != nil {
			return err
		}
		if !exists {
			if err := db.Exec(spec.definition).Error; err != nil {
				return err
			}
		}
		if spec.indexSQL != "" {
			if err := ensureIndexBySQL(db, "upload_files", spec.indexSQL); err != nil {
				return err
			}
		}
	}

	return nil
}

func hasColumn(db *gorm.DB, table, column string) (bool, error) {
	const countColumnSQL = `
SELECT COUNT(*)
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = ?
  AND COLUMN_NAME = ?
`

	var count int64
	if err := db.Raw(countColumnSQL, table, column).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func backfillUploadProvider(db *gorm.DB) error {
	type uploadPathRow struct {
		ID   int64
		Path string
	}

	var rows []uploadPathRow
	if err := db.Raw("SELECT id, path FROM upload_files WHERE provider = '' OR provider IS NULL OR provider = 'local'").Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		provider := detectUploadProvider(row.Path)
		if err := db.Exec("UPDATE upload_files SET provider = ? WHERE id = ?", provider, row.ID).Error; err != nil {
			return err
		}
	}

	return nil
}

func ensureIndexBySQL(db *gorm.DB, table, createSQL string) error {
	const countIndexSQL = `
SELECT COUNT(*)
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = ?
  AND INDEX_NAME = ?
`

	indexName := extractIndexName(createSQL)
	if indexName == "" {
		return nil
	}

	var count int64
	if err := db.Raw(countIndexSQL, table, indexName).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	return db.Exec(createSQL).Error
}

func extractIndexName(sql string) string {
	parts := strings.Fields(sql)
	for i, part := range parts {
		if strings.EqualFold(part, "INDEX") && i+1 < len(parts) {
			return strings.Trim(parts[i+1], "`")
		}
	}
	return ""
}

func detectUploadProvider(path string) string {
	switch {
	case strings.Contains(path, ".cos.") || strings.Contains(path, ".myqcloud.com"):
		return "tencent-cos"
	case strings.Contains(path, ".obs.") || strings.Contains(path, "huaweicloud"):
		return "huawei-obs"
	case strings.Contains(path, ".aliyuncs.com") || strings.Contains(path, ".oss-"):
		return "aliyun-oss"
	default:
		return "local"
	}
}
