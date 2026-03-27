package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/middleware"
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

func (h *ChannelHandler) Create(c *gin.Context) { h.save(c, 0) }

func (h *ChannelHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.save(c, id)
}

func (h *ChannelHandler) save(c *gin.Context, id int64) {
	var req dto.ChannelUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	claims := middleware.Claims(c)
	if err := h.service.Save(id, claims.AdminID, req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *ChannelHandler) BatchStatus(c *gin.Context) {
	var req dto.ChannelBatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.UpdateStatus(req.IDs, req.Status); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"updated": len(req.IDs)})
}
