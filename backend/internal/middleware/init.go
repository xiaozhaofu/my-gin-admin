package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func InitMiddleware(r *gin.Engine) {
	r.Use(Cors())
	r.Use(ErrorHandler())
	r.Use(requestid.New())
	r.Use(gin.Recovery())
}
