package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/middleware"
	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler { return &AuthHandler{service: service} }

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	data, err := h.service.Login(req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 4001, err.Error())
		return
	}
	response.OK(c, data)
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.Claims(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, 4001, "unauthorized")
		return
	}
	profile, err := h.service.Profile(claims.AdminID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, profile)
}
