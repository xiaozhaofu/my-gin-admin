package http

import (
	"go_sleep_admin/internal/handler"
	v2router "go_sleep_admin/internal/router"

	"github.com/gin-gonic/gin"
)

type Module struct {
	systemHandler        *handler.SystemHandler
	logHandler           *handler.LogHandler
	sessionHandler       *handler.SessionHandler
	jobHandler           *handler.JobHandler
	dataScopeHandler     *handler.DataScopeHandler
	authMiddleware       gin.HandlerFunc
	permissionMiddleware gin.HandlerFunc
	operationMiddleware  gin.HandlerFunc
}

func NewModule(
	systemHandler *handler.SystemHandler,
	logHandler *handler.LogHandler,
	sessionHandler *handler.SessionHandler,
	jobHandler *handler.JobHandler,
	dataScopeHandler *handler.DataScopeHandler,
	authMiddleware, permissionMiddleware, operationMiddleware gin.HandlerFunc,
) *Module {
	return &Module{
		systemHandler:        systemHandler,
		logHandler:           logHandler,
		sessionHandler:       sessionHandler,
		jobHandler:           jobHandler,
		dataScopeHandler:     dataScopeHandler,
		authMiddleware:       authMiddleware,
		permissionMiddleware: permissionMiddleware,
		operationMiddleware:  operationMiddleware,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("")
	useMiddlewares(group, m.operationMiddleware, m.authMiddleware, m.permissionMiddleware)

	group.GET("/dict-types", m.systemHandler.DictTypes)
	group.POST("/dict-types", m.systemHandler.SaveDictType)
	group.PUT("/dict-types/:id", m.systemHandler.SaveDictType)
	group.DELETE("/dict-types/:id", m.systemHandler.DeleteDictType)
	group.POST("/dict-items", m.systemHandler.SaveDictItem)
	group.PUT("/dict-items/:id", m.systemHandler.SaveDictItem)
	group.DELETE("/dict-items/:id", m.systemHandler.DeleteDictItem)

	group.GET("/sys-configs", m.systemHandler.Configs)
	group.POST("/sys-configs", m.systemHandler.SaveConfig)
	group.PUT("/sys-configs/:id", m.systemHandler.SaveConfig)
	group.DELETE("/sys-configs/:id", m.systemHandler.DeleteConfig)

	group.GET("/operation-logs", m.logHandler.List)
	group.DELETE("/operation-logs/:id", m.logHandler.Delete)
	group.DELETE("/operation-logs", m.logHandler.Clear)

	group.GET("/login-logs", m.sessionHandler.LoginLogs)
	group.GET("/online-sessions", m.sessionHandler.OnlineSessions)
	group.DELETE("/online-sessions/:id", m.sessionHandler.ForceOffline)

	group.GET("/jobs", m.jobHandler.List)
	group.POST("/jobs", m.jobHandler.Save)
	group.PUT("/jobs/:id", m.jobHandler.Save)
	group.DELETE("/jobs/:id", m.jobHandler.Delete)

	group.GET("/depts/tree", m.dataScopeHandler.DeptTree)
	group.POST("/depts", m.dataScopeHandler.SaveDept)
	group.PUT("/depts/:id", m.dataScopeHandler.SaveDept)
	group.DELETE("/depts/:id", m.dataScopeHandler.DeleteDept)

	group.GET("/posts", m.dataScopeHandler.Posts)
	group.POST("/posts", m.dataScopeHandler.SavePost)
	group.PUT("/posts/:id", m.dataScopeHandler.SavePost)
	group.DELETE("/posts/:id", m.dataScopeHandler.DeletePost)
}

func useMiddlewares(group *gin.RouterGroup, middlewares ...gin.HandlerFunc) {
	usable := make([]gin.HandlerFunc, 0, len(middlewares))
	for _, middleware := range middlewares {
		if middleware != nil {
			usable = append(usable, middleware)
		}
	}
	if len(usable) > 0 {
		group.Use(usable...)
	}
}

var _ v2router.Module = (*Module)(nil)
