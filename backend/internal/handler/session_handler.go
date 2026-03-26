package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/service"
)

type SessionHandler struct {
	service *service.SessionService
}

func NewSessionHandler(service *service.SessionService) *SessionHandler {
	return &SessionHandler{service: service}
}

func (h *SessionHandler) LoginLogs(c *gin.Context) {
	items, err := h.service.LoginLogs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *SessionHandler) OnlineSessions(c *gin.Context) {
	items, err := h.service.OnlineSessions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5000, err.Error())
		return
	}
	response.OK(c, items)
}

func (h *SessionHandler) ForceOffline(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.ForceOffline(id); err != nil {
		response.Error(c, http.StatusBadRequest, 4000, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}
