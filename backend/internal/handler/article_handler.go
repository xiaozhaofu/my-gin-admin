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

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(service *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) List(c *gin.Context) {
	var query dto.ArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	items, total, err := h.service.List(query, middleware.Claims(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	page, pageSize := query.Normalize()
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "page_size": pageSize})
}

func (h *ArticleHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.Detail(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 4004, err.Error())
		return
	}
	response.OK(c, item)
}

func (h *ArticleHandler) Create(c *gin.Context) { h.save(c, 0) }
func (h *ArticleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.save(c, id)
}

func (h *ArticleHandler) save(c *gin.Context, id int64) {
	var req dto.ArticleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	claims := middleware.Claims(c)
	if err := h.service.Save(id, claims.AdminID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *ArticleHandler) BatchStatus(c *gin.Context) {
	var req dto.ArticleBatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.UpdateStatus(req.IDs, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, gin.H{"updated": len(req.IDs)})
}

func (h *ArticleHandler) BatchCreate(c *gin.Context) {
	var req dto.ArticleBatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	claims := middleware.Claims(c)
	ids, err := h.service.BatchCreate(claims.AdminID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"created": len(ids), "ids": ids})
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	var req dto.BatchIDs
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.Delete(req.IDs); err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": len(req.IDs)})
}
