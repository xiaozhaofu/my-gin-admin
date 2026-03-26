package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, httpStatus int, code int, message string, data interface{}) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, 0, "success", data)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, 0, "success", data)
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Body{
		Code:    code,
		Message: message,
	})
}
