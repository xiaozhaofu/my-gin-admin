package http

import (
	"go_sleep_admin/internal/handler"
	v2router "go_sleep_admin/internal/router"

	"github.com/gin-gonic/gin"
)

type Module struct {
	menuHandler          *handler.MenuHandler
	articleHandler       *handler.ArticleHandler
	channelHandler       *handler.ChannelHandler
	uploadHandler        *handler.UploadHandler
	authMiddleware       gin.HandlerFunc
	permissionMiddleware gin.HandlerFunc
	operationMiddleware  gin.HandlerFunc
}

func NewModule(
	menuHandler *handler.MenuHandler,
	articleHandler *handler.ArticleHandler,
	channelHandler *handler.ChannelHandler,
	uploadHandler *handler.UploadHandler,
	authMiddleware, permissionMiddleware, operationMiddleware gin.HandlerFunc,
) *Module {
	return &Module{
		menuHandler:          menuHandler,
		articleHandler:       articleHandler,
		channelHandler:       channelHandler,
		uploadHandler:        uploadHandler,
		authMiddleware:       authMiddleware,
		permissionMiddleware: permissionMiddleware,
		operationMiddleware:  operationMiddleware,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("")
	useMiddlewares(group, m.operationMiddleware, m.authMiddleware, m.permissionMiddleware)

	group.GET("/menus/tree", m.menuHandler.Tree)
	group.GET("/menus/cascader", m.menuHandler.Cascader)
	group.GET("/channels", m.channelHandler.List)
	group.POST("/channels", m.channelHandler.Create)
	group.PUT("/channels/:id", m.channelHandler.Update)
	group.PUT("/channels/status", m.channelHandler.BatchStatus)
	group.POST("/menus", m.menuHandler.Create)
	group.PUT("/menus/:id", m.menuHandler.Update)
	group.DELETE("/menus/:id", m.menuHandler.Delete)

	group.GET("/articles", m.articleHandler.List)
	group.GET("/articles/:id", m.articleHandler.Detail)
	group.POST("/articles", m.articleHandler.Create)
	group.POST("/articles/batch", m.articleHandler.BatchCreate)
	group.PUT("/articles/:id", m.articleHandler.Update)
	group.PUT("/articles/status", m.articleHandler.BatchStatus)
	group.DELETE("/articles", m.articleHandler.Delete)

	group.POST("/uploads", m.uploadHandler.Upload)
	group.GET("/uploads", m.uploadHandler.List)
	group.DELETE("/uploads", m.uploadHandler.BatchDelete)
	group.DELETE("/uploads/:id", m.uploadHandler.Delete)
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
