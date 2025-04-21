package usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/DouglasValerio/cubiq-api/model"
	"github.com/DouglasValerio/cubiq-api/repository"
	"github.com/DouglasValerio/cubiq-api/service"
)

type InboundMessageUseCase struct {
	ChatMessageRepo    repository.ChatMessageRepository
	InboundMessageRepo repository.InboundMessageRepository
}

func (uc *InboundMessageUseCase) Execute(inboundMessage *model.InboundMessage) ([]*model.ChatMessage, error) {
	_, err := uc.InboundMessageRepo.InsertInboundMessage(inboundMessage)
	if err != nil {
		return nil, err
	}
	messages, statuses := inboundMessage.FlattenValues()
	chatMessages, errMsgs := uc.generateChatMessage(messages)
	chatStatuses, err := uc.generateChatStatus(statuses)

	for _, status := range chatStatuses {
		err := uc.ChatMessageRepo.UpdateChatMessageStatus(status.Id, derefStatuses(chatStatuses))
		if err != nil {
			return nil, err
		}
	}
	for _, message := range chatMessages {
		_, err := uc.ChatMessageRepo.InsertChatMessage(message)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	if errMsgs != nil {
		return nil, err
	}
	err = uc.InboundMessageRepo.SetProcessedAt(inboundMessage.ID)
	if err != nil {
		return nil, err
	}
	return chatMessages, nil
}
func (uc *InboundMessageUseCase) generateChatMessage(messages []model.Message) ([]*model.ChatMessage, error) {
	chatMessages := make([]*model.ChatMessage, 0)
	for _, message := range messages {
		t := time.Now().UTC()
		secondsSinceEpoch, err := strconv.ParseInt(message.Timestamp, 10, 64)
		if err == nil {

			t = time.Unix(secondsSinceEpoch, 0).UTC()
		}

		chatMessage := &model.ChatMessage{
			Id:        message.Id,
			From:      message.From,
			Type:      message.Type,
			Timestamp: t,
			Text:      message.Text,
		}
		chatMessage.SetDefaults()
		chatMessage.SetSummary()
		handleMediaMessage(message, chatMessage)
		chatMessages = append(chatMessages, chatMessage)
	}

	return chatMessages, nil
}

func handleMediaMessage(message model.Message, chatMessage *model.ChatMessage) {
	waService := service.WhatsAppService{}
	firebaseService := service.FirebaseUploader{}
	var mediaId string
	var contentType string

	if message.Type == "audio" {
		mediaId = message.Audio.Id
		contentType = message.Audio.MimeType
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Audio.Caption,
			Filename: message.Audio.Filename,
			Id:       message.Audio.Id,
			MimeType: message.Audio.MimeType,
		}
	}
	if message.Type == "video" {
		mediaId = message.Video.Id
		contentType = message.Video.MimeType
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Video.Caption,
			Filename: message.Video.Filename,
			Id:       message.Video.Id,
			MimeType: message.Video.MimeType,
		}
	}
	if message.Type == "image" {
		mediaId = message.Image.Id
		contentType = message.Image.MimeType
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Image.Caption,
			Filename: message.Image.Filename,
			Id:       message.Image.Id,
			MimeType: message.Image.MimeType,
		}
	}
	if message.Type == "document" {
		mediaId = message.Document.Id
		contentType = message.Document.MimeType
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Document.Caption,
			Filename: message.Document.Filename,
			Id:       message.Document.Id,
			MimeType: message.Document.MimeType,
		}
	}
	if message.Type == "sticker" {
		mediaId = message.Sticker.Id
		contentType = message.Sticker.MimeType
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Sticker.Caption,
			Filename: message.Sticker.Filename,
			Id:       message.Sticker.Id,
			MimeType: message.Sticker.MimeType,
		}
	}
	if mediaId != "" && contentType != "" {
		bytes, _ := waService.DownloadMedia(mediaId, "378510502012262", os.Getenv("WA_TOKEN"))

		url, err := firebaseService.UploadFile(context.Background(), bytes, mediaId, contentType)
		if err == nil {
			chatMessage.SetMediaFileUrl(url)
		}
	}

}
func (uc *InboundMessageUseCase) generateChatStatus(statuses []model.Statuses) ([]*model.MessageStatus, error) {
	chatStatuses := make([]*model.MessageStatus, 0)
	for _, status := range statuses {
		chatMessage := &model.MessageStatus{
			Status:    status.Status,
			Timestamp: status.Timestamp,
			Id:        status.Id,
		}
		chatStatuses = append(chatStatuses, chatMessage)
	}

	return chatStatuses, nil
}

func derefStatuses(ptrs []*model.MessageStatus) []model.MessageStatus {
	vals := make([]model.MessageStatus, len(ptrs))
	for i, p := range ptrs {
		if p != nil {
			vals[i] = *p
		}
	}
	return vals
}
