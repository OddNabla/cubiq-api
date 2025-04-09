package model

import (
	"time"

	"github.com/google/uuid"
)

type InboundMessage struct {
	ID        string                `json:"id,omitempty" bson:"_id"`
	Object    string                `json:"object" binding:"required" bson:"object"`
	Entry     []InboundMessageEntry `json:"entry" binding:"required" bson:"entry"`
	CreatedAt time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt" bson:"updatedAt"`
}

type InboundMessageEntry struct {
	Id      string                       `json:"id" binding:"required" bson:"id"`
	Changes []InboundMessageEntryChanges `json:"changes" binding:"required" bson:"changes"`
}

type InboundMessageEntryChanges struct {
	Value InboundMessageEntryValue `json:"value" binding:"required" bson:"value"`
	Field string                   `json:"field" binding:"required" bson:"field"`
}
type InboundMessageEntryValue struct {
	MessagingProduct string     `json:"messaging_product" binding:"required" bson:"messaging_product"`
	Metadata         Metadata   `json:"metadata" binding:"required" bson:"metadata"`
	Contacts         []Contact  `json:"contacts" binding:"required" bson:"contacts"`
	Messages         []Message  `json:"messages" binding:"required" bson:"messages"`
	Statuses         []Statuses `json:"statuses" binding:"required" bson:"statuses"`
}
type Metadata struct {
	PhoneNumberId      string `json:"phone_number_id" binding:"required" bson:"phone_number_id"`
	DisplayPhoneNumber string `json:"display_phone_number" binding:"required" bson:"display_phone_number"`
}

type Contact struct {
	Profile ContactProfile `json:"profile" binding:"required" bson:"profile"`
	WaId    string         `json:"wa_id" binding:"required" bson:"wa_id"`
}

type ContactProfile struct {
	Name string `json:"name" binding:"required" bson:"name"`
}
type Message struct {
	From      string      `json:"from" binding:"required" bson:"from"`
	Id        string      `json:"id" binding:"required" bson:"id"`
	Timestamp string      `json:"timestamp" binding:"required" bson:"timestamp"`
	Type      string      `json:"type" binding:"required" bson:"type"`
	Text      MessageText `json:"text" binding:"required" bson:"text"`
}
type MessageText struct {
	Body       string `json:"body" binding:"required" bson:"body"`
	PreviewUrl bool   `json:"preview_url" binding:"required" bson:"preview_url"`
}

type Statuses struct {
	Id           string       `json:"id" binding:"required" bson:"id"`
	Status       string       `json:"status" binding:"required" bson:"status"`
	Timestamp    string       `json:"timestamp" binding:"required" bson:"timestamp"`
	RecipientId  string       `json:"recipient_id" binding:"required" bson:"recipient_id"`
	Conversation Conversation `json:"conversation" binding:"required" bson:"conversation"`
	Pricing      Pricing      `json:"pricing" binding:"required" bson:"pricing"`
}

type Conversation struct {
	Id     string `json:"id" binding:"required" bson:"id"`
	Origin Origin `json:"origin" binding:"required" bson:"origin"`
}

type Origin struct {
	Type string `json:"type" binding:"required" bson:"type"`
}

type Pricing struct {
	Category     string `json:"category" binding:"required" bson:"category"`
	Billable     bool   `json:"billable" binding:"required" bson:"billable"`
	PricingModel string `json:"pricing_model" binding:"required" bson:"pricing_model"`
}

func (m *InboundMessage) SetDefaults() {
	now := time.Now().UTC()
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
}
