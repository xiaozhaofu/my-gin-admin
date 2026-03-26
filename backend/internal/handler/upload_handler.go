package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/middleware"
	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/platform/storage"
	"go_sleep_admin/internal/service"
)

type UploadHandler struct {
	service  *service.UploadService
	uploader storage.UploadGateway
}

func NewUploadHandler(service *service.UploadService, uploader storage.UploadGateway) *UploadHandler {
	return &UploadHandler{service: service, uploader: uploader}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	scene := c.DefaultPostForm("scene", "misc")
	provider := c.DefaultPostForm("provider", "")
	saved, err := h.uploader.Save(file, scene, provider)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	claims := middleware.Claims(c)
	if err := h.service.Upload(claims.AdminID, saved); err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, gin.H{
		"id":       saved.Random,
		"url":      saved.Path,
		"name":     saved.OriginName,
		"ftype":    saved.Type,
		"md5":      saved.MD5,
		"scene":    saved.Scene,
		"provider": saved.Provider,
	})
}

func (h *UploadHandler) List(c *gin.Context) {
	var query dto.UploadListQuery
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

func (h *UploadHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *UploadHandler) BatchDelete(c *gin.Context) {
	var req dto.BatchIDs
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.DeleteBatch(req.IDs); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": len(req.IDs)})
}
