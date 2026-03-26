package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type DashboardHandler struct {
	service *service.DashboardService
}

func NewDashboardHandler(service *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	data, err := h.service.Overview()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, data)
}
