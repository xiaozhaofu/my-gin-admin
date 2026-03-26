package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/middleware"
	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type RBACHandler struct {
	service *service.RBACService
}

func NewRBACHandler(service *service.RBACService) *RBACHandler { return &RBACHandler{service: service} }

func (h *RBACHandler) Roles(c *gin.Context) {
	items, err := h.service.Roles()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *RBACHandler) SaveRole(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.RoleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.SaveRole(id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *RBACHandler) Admins(c *gin.Context) {
	items, err := h.service.Admins(middleware.Claims(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *RBACHandler) SaveAdmin(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.AdminUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	passwordHash := ""
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, 5000, err.Error())
			return
		}
		passwordHash = string(hashed)
	}
	if err := h.service.SaveAdmin(id, passwordHash, req); err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *RBACHandler) AdminMenus(c *gin.Context) {
	items, err := h.service.AdminMenuTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *RBACHandler) SaveAdminMenu(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.AdminMenuUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	if err := h.service.SaveAdminMenu(id, req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *RBACHandler) DeleteAdminMenu(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteAdminMenu(id); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}
