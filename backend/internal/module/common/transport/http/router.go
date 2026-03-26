package http

import (
	"net/http"

	"go_sleep_admin/internal/handler"
	v2router "go_sleep_admin/internal/router"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler             *handler.ModuleHandler
	dashboardHandler    *handler.DashboardHandler
	authMiddleware      gin.HandlerFunc
	operationMiddleware gin.HandlerFunc
}

func NewModule(handler *handler.ModuleHandler, dashboardHandler *handler.DashboardHandler, authMiddleware, operationMiddleware gin.HandlerFunc) *Module {
	return &Module{
		handler:             handler,
		dashboardHandler:    dashboardHandler,
		authMiddleware:      authMiddleware,
		operationMiddleware: operationMiddleware,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
	})

	protected := router.Group("")
	useMiddlewares(protected, m.operationMiddleware, m.authMiddleware)
	protected.GET("/module-metadata", m.handler.List)
	protected.GET("/dashboard", m.dashboardHandler.Overview)
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
