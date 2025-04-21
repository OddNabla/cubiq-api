package handler

import (
	// "encoding/json"
	// "fmt"
	"net/http"

	"github.com/DouglasValerio/cubiq-api/model"
	"github.com/DouglasValerio/cubiq-api/repository"
	"github.com/DouglasValerio/cubiq-api/setup"
	"github.com/DouglasValerio/cubiq-api/usecase"
	"github.com/gin-gonic/gin"
)

func HandleInboundMessage(c *gin.Context) {
	var inboundMessage model.InboundMessage
	if err := c.ShouldBindJSON(&inboundMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inboundMessageRepository := repository.InboundMessageRepo{
		MongoDatabase: setup.MongoDatabase,
	}
	chatRepository := &repository.ChatMessageRepo{
		MongoDatabase: setup.MongoDatabase,
	}
	inboundMessageUseCase := usecase.InboundMessageUseCase{
		ChatMessageRepo: chatRepository}
	inboundMessage.SetDefaults()
	result, err := inboundMessageRepository.InsertInboundMessage(&inboundMessage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save inbound message"})
		return
	}
	_, err = inboundMessageUseCase.Execute(&inboundMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save inbound message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "result": result})
}
func HandleInboundFromWebhookMessage(c *gin.Context) {
	var inboundMessage model.InboundMessage
	if err := c.ShouldBindJSON(&inboundMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inboundMessageRepository := &repository.InboundMessageRepo{
		MongoDatabase: setup.MongoDatabase,
	}
	chatRepository := &repository.ChatMessageRepo{
		MongoDatabase: setup.MongoDatabase,
	}
	inboundMessageUseCase := usecase.InboundMessageUseCase{
		ChatMessageRepo:    chatRepository,
		InboundMessageRepo: inboundMessageRepository,
	}
	inboundMessage.SetDefaults()
	result, err := inboundMessageUseCase.Execute(&inboundMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save inbound message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "result": result})
}
