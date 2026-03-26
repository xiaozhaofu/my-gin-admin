package http

import (
	"go_sleep_admin/internal/handler"
	v2router "go_sleep_admin/internal/router"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler             *handler.AuthHandler
	authMiddleware      gin.HandlerFunc
	operationMiddleware gin.HandlerFunc
}

func NewModule(handler *handler.AuthHandler, authMiddleware, operationMiddleware gin.HandlerFunc) *Module {
	return &Module{
		handler:             handler,
		authMiddleware:      authMiddleware,
		operationMiddleware: operationMiddleware,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	useMiddlewares(authGroup, m.operationMiddleware)
	authGroup.POST("/login", m.handler.Login)

	protected := router.Group("/auth")
	useMiddlewares(protected, m.operationMiddleware, m.authMiddleware)
	protected.GET("/me", m.handler.Me)
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
