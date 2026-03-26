package handler

import (
	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type ModuleHandler struct {
	service *service.ModuleService
}

func NewModuleHandler(service *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{service: service}
}

func (h *ModuleHandler) List(c *gin.Context) {
	response.OK(c, h.service.Modules())
}
