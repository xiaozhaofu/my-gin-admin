package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtkit/logger"
	"go.uber.org/zap"

	"go_sleep_admin/internal/pkg/apperror"
	"go_sleep_admin/internal/pkg/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.HTTP, appErr.Code, appErr.Message)
			return
		}
		logger.ZError("request failed", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, 5000, "internal server error")
	}
}
