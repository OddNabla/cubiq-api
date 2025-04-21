package usecase_test

import (
	"testing"
	"time"

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

type MockInboundMessageRepo struct {
	mock.Mock
}

func (m *MockInboundMessageRepo) InsertInboundMessage(inbound *model.InboundMessage) (*model.InboundMessage, error) {
	return &model.InboundMessage{
		ID: inbound.ID,
	}, nil
}
func (m *MockInboundMessageRepo) SetProcessedAt(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockInboundMessageRepo) FindInboundMessageById(id string) (*model.InboundMessage, error) {
	args := m.Called(id)
	return args.Get(0).(*model.InboundMessage), args.Error(1)
}
func (m *MockInboundMessageRepo) FindAllInboundMessages() ([]model.InboundMessage, error) {
	args := m.Called()
	return args.Get(0).([]model.InboundMessage), args.Error(1)
}

func TestInboundMessageUseCase_Execute(t *testing.T) {
	mockRepo := &MockChatMessageRepo{}
	MockInboundMessageRepo := &MockInboundMessageRepo{}
	useCase := &usecase.InboundMessageUseCase{
		ChatMessageRepo:    mockRepo,
		InboundMessageRepo: MockInboundMessageRepo,
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
									Timestamp: "1745240279",
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

	parsedTime, err := time.Parse(time.UnixDate, time.Unix(1745240279, 0).UTC().Format(time.UnixDate))
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	expectedMessage := &model.ChatMessage{
		Id:        "msg-1",
		From:      "user-123",
		Timestamp: parsedTime,
		Type:      "text",
		Text:      model.MessageText{Body: "Hello"},
	}

	// Expectations
	mockRepo.On("InsertChatMessage", expectedMessage).Return(nil)

	MockInboundMessageRepo.On("InsertInboundMessage", inbound).Return(inbound, nil)
	MockInboundMessageRepo.On("SetProcessedAt", inbound.ID).Return(nil)

	// Run
	result, err := useCase.Execute(inbound)

	// Verify
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedMessage.Text.Body, result[0].Text.Body)
}
