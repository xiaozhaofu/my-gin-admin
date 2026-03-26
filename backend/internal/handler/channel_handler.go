package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type ChannelHandler struct {
	service *service.ChannelService
}

func NewChannelHandler(service *service.ChannelService) *ChannelHandler {
	return &ChannelHandler{service: service}
}

func (h *ChannelHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}
