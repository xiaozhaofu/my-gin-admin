package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/pkg/response"
	"go_sleep_admin/internal/platform/auth"
	"go_sleep_admin/internal/service"
)

const ContextAdminClaims = "admin_claims"

func JWT(jwtManager *auth.JWTManager, sessionSvc *service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer"))
		if token == "" {
			response.Error(c, http.StatusUnauthorized, 4001, "missing token")
			return
		}
		claims, err := jwtManager.Parse(token)
		if err != nil || claims.TokenUse != "access" {
			response.Error(c, http.StatusUnauthorized, 4001, "invalid token")
			return
		}
		if forced, err := sessionSvc.IsForcedOffline(token); err == nil && forced {
			response.Error(c, http.StatusUnauthorized, 4001, "session forced offline")
			return
		}
		_ = sessionSvc.Touch(token)
		c.Set("access_token", token)
		c.Set("request_time", time.Now())
		c.Set(ContextAdminClaims, claims)
		c.Next()
	}
}

func Claims(c *gin.Context) *auth.Claims {
	val, _ := c.Get(ContextAdminClaims)
	claims, _ := val.(*auth.Claims)
	return claims
}
