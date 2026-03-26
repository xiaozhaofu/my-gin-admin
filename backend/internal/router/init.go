package router

import (
	"strings"

	"go_sleep_admin/internal/appmeta"
	legacymiddleware "go_sleep_admin/internal/middleware"
	"go_sleep_admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// InitRouterWithModules 创建 Gin engine 并注册所有业务模块路由。
//
// 新增模块只需在 providers.go 中将 Module 实现追加到切片，
// 路由层无需任何改动。
func InitRouterWithModules(envName string, modules []Module) *gin.Engine {
	engine := gin.New()
	gin.SetMode(ginModeFromEnv(envName))
	engine.RedirectTrailingSlash = false
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(noMethod)
	engine.NoRoute(noRoute)
	legacymiddleware.InitMiddleware(engine)

	group := engine.Group(appmeta.APIV1Path)
	for _, m := range modules {
		m.RegisterRoutes(group)
	}

	return engine
}

func ginModeFromEnv(envName string) string {
	switch strings.ToLower(strings.TrimSpace(envName)) {
	case "prod":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}

func noRoute(c *gin.Context) {
	if strings.Contains(c.Request.URL.Path, ".") {
		return
	}
	response.Error(c, 404, 4004, "unknown route")
}

func noMethod(c *gin.Context) {
	response.Error(c, 405, 4005, "method not allowed")
}
