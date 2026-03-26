package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/service"
)

func OperationLogger(logService *service.LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.URL.Path == "/api/v1/healthz" {
			c.Next()
			return
		}

		start := time.Now()
		bodyText := captureRequestBody(c)
		c.Next()

		adminID := int64(0)
		username := ""
		if claims := Claims(c); claims != nil {
			adminID = claims.AdminID
			username = claims.Username
		}

		statusCode := c.Writer.Status()
		errMsg := ""
		if len(c.Errors) > 0 {
			errMsg = c.Errors.Last().Error()
		}

		_ = logService.Create(&models.OperationLog{
			AdminID:      adminID,
			Username:     username,
			Method:       c.Request.Method,
			Path:         c.FullPath(),
			StatusCode:   statusCode,
			Success:      statusCode < 400,
			ClientIP:     c.ClientIP(),
			UserAgent:    truncateLog(c.Request.UserAgent(), 500),
			RequestBody:  truncateLog(bodyText, 4000),
			DurationMS:   time.Since(start).Milliseconds(),
			ErrorMessage: truncateLog(errMsg, 1000),
		})
	}
}

func captureRequestBody(c *gin.Context) string {
	if c.Request == nil || c.Request.Body == nil {
		return ""
	}
	if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
		return "[multipart/form-data omitted]"
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

func truncateLog(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
