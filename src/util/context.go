package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data interface{}){
	c.JSON(http.StatusOK, data)
}

func Error(c *gin.Context, status int, msg string) {
	c.String(status, msg)
}