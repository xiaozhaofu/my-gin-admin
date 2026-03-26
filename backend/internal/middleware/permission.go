package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"go_sleep_admin/internal/pkg/response"
)

func Casbin(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := Claims(c)
		if claims == nil {
			response.Error(c, http.StatusUnauthorized, 4001, "unauthorized")
			return
		}

		for _, role := range claims.Roles {
			if role == "admin" {
				c.Next()
				return
			}
		}

		allowed := false
		for _, role := range claims.Roles {
			ok, err := enforcer.Enforce(role, c.FullPath(), c.Request.Method)
			if err == nil && ok {
				allowed = true
				break
			}
		}
		if !allowed {
			response.Error(c, http.StatusForbidden, 4003, "forbidden")
			return
		}
		c.Next()
	}
}
