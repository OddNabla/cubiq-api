package usecase

import (
	"github.com/DouglasValerio/cubiq-api/model"
	"github.com/DouglasValerio/cubiq-api/repository"
)

type InboundMessageUseCase struct {
	ChatMessageRepo repository.ChatMessageRepo
}

func (uc *InboundMessageUseCase) Execute(inboundMessage *model.InboundMessage) ([]*model.ChatMessage, error) {

	messages, statuses := inboundMessage.FlattenValues()
	chatMessages, errMsgs := uc.generateChatMessage(messages)
	chatStatuses, err := uc.generateChatStatus(statuses)

	for _, status := range chatStatuses {
		uc.ChatMessageRepo.UpdateChatMessageStatus(status.Id, derefStatuses(chatStatuses))
	}
	for _, message := range chatMessages {
		uc.ChatMessageRepo.InsertChatMessage(message)
	}

	if err != nil {
		return nil, err
	}

	if errMsgs != nil {
		return nil, err
	}

	return chatMessages, nil
}
func (uc *InboundMessageUseCase) generateChatMessage(messages []model.Message) ([]*model.ChatMessage, error) {
	chatMessages := make([]*model.ChatMessage, 0)
	for _, message := range messages {
		chatMessage := &model.ChatMessage{
			Id:        message.Id,
			From:      message.From,
			Type:      message.Type,
			Timestamp: message.Timestamp,
			Text:      message.Text,
		}
		handleMediaMessage(message, chatMessage)
		chatMessages = append(chatMessages, chatMessage)
	}

	return chatMessages, nil
}

func handleMediaMessage(message model.Message, chatMessage *model.ChatMessage) {
	if message.Type == "audio" {
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Audio.Caption,
			Filename: message.Audio.Filename,
			Id:       message.Audio.Id,
			MimeType: message.Audio.MimeType,
		}
	}
	if message.Type == "video" {
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Video.Caption,
			Filename: message.Video.Filename,
			Id:       message.Video.Id,
			MimeType: message.Video.MimeType,
		}
	}
	if message.Type == "image" {
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Image.Caption,
			Filename: message.Image.Filename,
			Id:       message.Image.Id,
			MimeType: message.Image.MimeType,
		}
	}
	if message.Type == "document" {
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Document.Caption,
			Filename: message.Document.Filename,
			Id:       message.Document.Id,
			MimeType: message.Document.MimeType,
		}
	}
	if message.Type == "sticker" {
		chatMessage.Media = model.MediaMessage{
			Caption:  message.Sticker.Caption,
			Filename: message.Sticker.Filename,
			Id:       message.Sticker.Id,
			MimeType: message.Sticker.MimeType,
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
