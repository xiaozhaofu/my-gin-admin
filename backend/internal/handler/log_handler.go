package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type LogHandler struct {
	service *service.LogService
}

func NewLogHandler(service *service.LogService) *LogHandler { return &LogHandler{service: service} }

func (h *LogHandler) List(c *gin.Context) {
	var query dto.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	items, total, err := h.service.List(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	page, pageSize := query.Normalize()
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "page_size": pageSize})
}

func (h *LogHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *LogHandler) Clear(c *gin.Context) {
	if err := h.service.Clear(); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"cleared": true})
}
