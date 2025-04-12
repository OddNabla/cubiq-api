package usecase_test

import (
	"testing"

	"github.com/DouglasValerio/cubiq-api/model"
	"github.com/DouglasValerio/cubiq-api/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChatMessageRepo struct {
	mock.Mock
}

func (m *MockChatMessageRepo) InsertChatMessage(msg *model.ChatMessage) (*model.ChatMessage, error) {
	return &model.ChatMessage{
		Id:        msg.Id,
		From:      msg.From,
		Timestamp: msg.Timestamp,
		Type:      msg.Type,
		Text:      msg.Text,
		Statuses:  msg.Statuses,
		Media:     msg.Media,
	}, nil
}

func (m *MockChatMessageRepo) UpdateChatMessageStatus(id string, statuses []model.MessageStatus) error {
	args := m.Called(id, statuses)
	return args.Error(0)
}

func TestInboundMessageUseCase_Execute(t *testing.T) {
	mockRepo := new(MockChatMessageRepo)
	useCase := &usecase.InboundMessageUseCase{
		ChatMessageRepo: mockRepo,
	}

	inbound := &model.InboundMessage{
		Entry: []model.InboundMessageEntry{
			{
				Id: "msg-1",
				Changes: []model.InboundMessageEntryChanges{
					{
						Value: model.InboundMessageEntryValue{
							MessagingProduct: "whatsapp",
							Metadata: model.Metadata{
								PhoneNumberId:      "1234567890",
								DisplayPhoneNumber: "1234567890",
							},
							Messages: []model.Message{
								{
									Id:        "msg-1",
									From:      "user-123",
									Timestamp: "1234567890",
									Type:      "text",
									Text:      model.MessageText{Body: "Hello"},
								},
							},
						},
					},
				},
			},
		},
	}

	expectedStatus := model.MessageStatus{
		Id:        "msg-1",
		Status:    "delivered",
		Timestamp: "1234567891",
	}

	expectedMessage := &model.ChatMessage{
		Id:        "msg-1",
		From:      "user-123",
		Timestamp: "1234567890",
		Type:      "text",
		Text:      model.MessageText{Body: "Hello"},
	}

	// Expectations
	mockRepo.On("UpdateChatMessageStatus", "msg-1", []model.MessageStatus{expectedStatus}).Return(nil)
	mockRepo.On("InsertChatMessage", expectedMessage).Return(nil)

	// Run
	result, err := useCase.Execute(inbound)

	// Verify
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedMessage.Text.Body, result[0].Text.Body)

	mockRepo.AssertExpectations(t)
}
