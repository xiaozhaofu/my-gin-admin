package http

import (
	"go_sleep_admin/internal/handler"
	v2router "go_sleep_admin/internal/router"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler              *handler.RBACHandler
	authMiddleware       gin.HandlerFunc
	permissionMiddleware gin.HandlerFunc
	operationMiddleware  gin.HandlerFunc
}

func NewModule(handler *handler.RBACHandler, authMiddleware, permissionMiddleware, operationMiddleware gin.HandlerFunc) *Module {
	return &Module{
		handler:              handler,
		authMiddleware:       authMiddleware,
		permissionMiddleware: permissionMiddleware,
		operationMiddleware:  operationMiddleware,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("")
	useMiddlewares(group, m.operationMiddleware, m.authMiddleware, m.permissionMiddleware)

	group.GET("/admins", m.handler.Admins)
	group.POST("/admins", m.handler.SaveAdmin)
	group.PUT("/admins/:id", m.handler.SaveAdmin)

	group.GET("/roles", m.handler.Roles)
	group.POST("/roles", m.handler.SaveRole)
	group.PUT("/roles/:id", m.handler.SaveRole)

	group.GET("/admin-menus/tree", m.handler.AdminMenus)
	group.POST("/admin-menus", m.handler.SaveAdminMenu)
	group.PUT("/admin-menus/:id", m.handler.SaveAdminMenu)
	group.DELETE("/admin-menus/:id", m.handler.DeleteAdminMenu)
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
