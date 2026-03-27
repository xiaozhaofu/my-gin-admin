package bootstrap

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go_sleep_admin/internal/handler"
	"go_sleep_admin/internal/platform/auth"
	"go_sleep_admin/internal/platform/storage"
	"go_sleep_admin/internal/repository"
	"go_sleep_admin/internal/service"
)

type repositories struct {
	Admin     *repository.AdminRepository
	AdminMenu *repository.AdminMenuRepository
	Article   *repository.ArticleRepository
	Channel   *repository.ChannelRepository
	DataScope *repository.DataScopeRepository
	Dashboard *repository.DashboardRepository
	Job       *repository.JobRepository
	Log       *repository.LogRepository
	Menu      *repository.MenuRepository
	Order     *repository.OrderRepository
	Role      *repository.RoleRepository
	Session   *repository.SessionRepository
	System    *repository.SystemRepository
	Upload    *repository.UploadRepository
}

func newRepositories(db *gorm.DB) repositories {
	return repositories{
		Admin:     repository.NewAdminRepository(db),
		AdminMenu: repository.NewAdminMenuRepository(db),
		Article:   repository.NewArticleRepository(db),
		Channel:   repository.NewChannelRepository(db),
		DataScope: repository.NewDataScopeRepository(db),
		Dashboard: repository.NewDashboardRepository(db),
		Job:       repository.NewJobRepository(db),
		Log:       repository.NewLogRepository(db),
		Menu:      repository.NewMenuRepository(db),
		Order:     repository.NewOrderRepository(db),
		Role:      repository.NewRoleRepository(db),
		Session:   repository.NewSessionRepository(db),
		System:    repository.NewSystemRepository(db),
		Upload:    repository.NewUploadRepository(db),
	}
}

type services struct {
	Article   *service.ArticleService
	Auth      *service.AuthService
	Channel   *service.ChannelService
	DataScope *service.DataScopeService
	Dashboard *service.DashboardService
	Job       *service.JobService
	Log       *service.LogService
	Menu      *service.MenuService
	Module    *service.ModuleService
	Order     *service.OrderService
	RBAC      *service.RBACService
	Session   *service.SessionService
	System    *service.SystemService
	Upload    *service.UploadService
}

func newServices(db *gorm.DB, enforcer *casbin.Enforcer, jwtManager *auth.JWTManager, repos repositories) services {
	sessionService := service.NewSessionService(repos.Session)

	return services{
		Article:   service.NewArticleService(repos.Article, repos.Admin, repos.DataScope),
		Auth:      service.NewAuthService(repos.Admin, enforcer, jwtManager, sessionService),
		Channel:   service.NewChannelService(repos.Channel),
		DataScope: service.NewDataScopeService(repos.DataScope),
		Dashboard: service.NewDashboardService(repos.Dashboard),
		Job:       service.NewJobService(repos.Job),
		Log:       service.NewLogService(repos.Log),
		Menu:      service.NewMenuService(repos.Menu),
		Module:    service.NewModuleService(),
		Order:     service.NewOrderService(repos.Order),
		RBAC:      service.NewRBACService(db, repos.Role, repos.Admin, repos.AdminMenu, repos.DataScope, enforcer),
		Session:   sessionService,
		System:    service.NewSystemService(repos.System),
		Upload:    service.NewUploadService(repos.Upload, repos.Admin, repos.DataScope),
	}
}

type handlers struct {
	Article   *handler.ArticleHandler
	Auth      *handler.AuthHandler
	Channel   *handler.ChannelHandler
	DataScope *handler.DataScopeHandler
	Dashboard *handler.DashboardHandler
	Job       *handler.JobHandler
	Log       *handler.LogHandler
	Menu      *handler.MenuHandler
	Module    *handler.ModuleHandler
	Order     *handler.OrderHandler
	RBAC      *handler.RBACHandler
	Session   *handler.SessionHandler
	System    *handler.SystemHandler
	Upload    *handler.UploadHandler
}

func newHandlers(services services, uploader storage.UploadGateway) handlers {
	return handlers{
		Article:   handler.NewArticleHandler(services.Article),
		Auth:      handler.NewAuthHandler(services.Auth),
		Channel:   handler.NewChannelHandler(services.Channel),
		DataScope: handler.NewDataScopeHandler(services.DataScope),
		Dashboard: handler.NewDashboardHandler(services.Dashboard),
		Job:       handler.NewJobHandler(services.Job),
		Log:       handler.NewLogHandler(services.Log),
		Menu:      handler.NewMenuHandler(services.Menu),
		Module:    handler.NewModuleHandler(services.Module),
		Order:     handler.NewOrderHandler(services.Order),
		RBAC:      handler.NewRBACHandler(services.RBAC),
		Session:   handler.NewSessionHandler(services.Session),
		System:    handler.NewSystemHandler(services.System),
		Upload:    handler.NewUploadHandler(services.Upload, uploader),
	}
}

type httpComponents struct {
	Services services
	Handlers handlers
}

func newHTTPComponents(db *gorm.DB, enforcer *casbin.Enforcer, jwtManager *auth.JWTManager, uploader storage.UploadGateway) httpComponents {
	repos := newRepositories(db)
	services := newServices(db, enforcer, jwtManager, repos)
	return httpComponents{
		Services: services,
		Handlers: newHandlers(services, uploader),
	}
}
