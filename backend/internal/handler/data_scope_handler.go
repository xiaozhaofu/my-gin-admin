package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type DataScopeHandler struct {
	service *service.DataScopeService
}

func NewDataScopeHandler(service *service.DataScopeService) *DataScopeHandler {
	return &DataScopeHandler{service: service}
}

func (h *DataScopeHandler) DeptTree(c *gin.Context) {
	items, err := h.service.DeptTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *DataScopeHandler) SaveDept(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.DeptUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.SaveDept(id, req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *DataScopeHandler) DeleteDept(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteDept(id); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *DataScopeHandler) Posts(c *gin.Context) {
	items, err := h.service.Posts()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *DataScopeHandler) SavePost(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.PostUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.SavePost(id, req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *DataScopeHandler) DeletePost(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeletePost(id); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}
