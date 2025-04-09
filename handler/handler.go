package handler

import (
	"net/http"

	"github.com/DouglasValerio/cubiq-api/model"
	"github.com/DouglasValerio/cubiq-api/repository"
	"github.com/DouglasValerio/cubiq-api/setup"
	"github.com/gin-gonic/gin"
)

func HandleInboundMessage(c *gin.Context) {
	var inboundMessage model.InboundMessage
	if err := c.ShouldBindJSON(&inboundMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repository := repository.InboundMessageRepo{
		MongoCollection: setup.MongoClient.Database("cubiq").Collection("inbound_messages"),
	}
	result, err := repository.InsertInboundMessage(&inboundMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save inbound message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "result": result})
}
