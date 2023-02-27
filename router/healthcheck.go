package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadHealthCheck(e *gin.Engine) {
	e.GET("health-check", healthCheck)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
