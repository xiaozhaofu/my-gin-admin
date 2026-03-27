package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/repository"
	"go_sleep_admin/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) List(c *gin.Context) {
	var query dto.OrderQuery
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
	response.OK(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *OrderHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.Detail(id)
	if err != nil {
		if repository.IsOrderNotFound(err) {
			response.Error(c, http.StatusNotFound, 4004, "订单不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, item)
}

func (h *OrderHandler) Export(c *gin.Context) {
	var query dto.OrderQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}

	data, err := h.service.Export(query)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}

	filename := fmt.Sprintf("orders_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
