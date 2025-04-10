package router

import (
	"net/http"

	"github.com/DouglasValerio/cubiq-api/handler"
	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		v1.POST("/inbound-message", handler.HandleInboundMessage)
	}
	v2 := router.Group("/api/chat")
	{
		v2.POST("/wa", handler.HandleInboundFromWebhookMessage)
	}
}
