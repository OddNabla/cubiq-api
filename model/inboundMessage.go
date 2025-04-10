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
	From      string           `json:"from" binding:"required" bson:"from"`
	Id        string           `json:"id" binding:"required" bson:"id"`
	Timestamp string           `json:"timestamp" binding:"required" bson:"timestamp"`
	Type      string           `json:"type" binding:"required" bson:"type"`
	Text      MessageText      `json:"text" bson:"text"`
	Context   Context          `json:"context" bson:"context"`
	Audio     Media            `json:"audio" bson:"audio"`
	Image     Media            `json:"image" bson:"image"`
	Video     Media            `json:"video" bson:"video"`
	Document  Media            `json:"document" bson:"document"`
	Sticker   Media            `json:"sticker" bson:"sticker"`
	Location  Location         `json:"location" bson:"location"`
	Contacts  []SharedContacts `json:"contacts" bson:"contacts"`
	Reaction  Reaction         `json:"reaction" bson:"reaction"`
	Button    Button           `json:"button" bson:"button"`
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

type Context struct {
	From                string `json:"from" bson:"from"`
	Id                  string `json:"id" bson:"id"`
	Forwarded           bool   `json:"forwarded" bson:"forwarded"`
	FrequentlyForwarded bool   `json:"frequently_forwarded" bson:"frequently_forwarded"`
	ReferredProduct     string `json:"referred_product" bson:"referred_product"`
}

type Media struct {
	Caption  string `json:"caption" bson:"caption"`
	Filename string `json:"filename" bson:"filename"`
	Id       string `json:"id" binding:"required" bson:"id"`
	MimeType string `json:"mime_type" binding:"required" bson:"mime_type"`
	Sha256   string `json:"sha256" binding:"required" bson:"sha256"`
}

type SharedContacts struct {
	Name   SharedContactName    `json:"name" binding:"required" bson:"name"`
	Phones []SharedContactPhone `json:"phones" binding:"required" bson:"phones"`
}

type SharedContactName struct {
	FirstName     string `json:"first_name" binding:"required" bson:"first_name"`
	FormattedName string `json:"formatted_name" binding:"required" bson:"formatted_name"`
	LastName      string `json:"last_name" binding:"required" bson:"last_name"`
	MiddleName    string `json:"middle_name" binding:"required" bson:"middle_name"`
}

type SharedContactPhone struct {
	Phone string `json:"phone" binding:"required" bson:"phone"`
	Type  string `json:"type" binding:"required" bson:"type"`
	WaId  string `json:"wa_id" binding:"required" bson:"wa_id"`
}

type Reaction struct {
	MessageId string `json:"message_id" binding:"required" bson:"message_id"`
	Emoji     string `json:"emoji" binding:"required" bson:"emoji"`
}
type Location struct {
	Address   string  `json:"address" bson:"address"`
	Latitude  float32 `json:"latitude" bson:"latitude"`
	Longitude float32 `json:"longitude" bson:"longitude"`
	Name      string  `json:"name" bson:"name"`
}

type Button struct {
	Text    string `json:"text" bson:"text"`
	Payload string `json:"payload" bson:"payload"`
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

func (m *InboundMessage) FlattenValues() ([]Message, []Statuses) {
	var allMessages []Message
	var allStatuses []Statuses
	for _, entry := range m.Entry {
		for _, change := range entry.Changes {
			if change.Value.Messages != nil {
				allMessages = append(allMessages, change.Value.Messages...)
			}
		}
	}
	for _, entry := range m.Entry {
		for _, change := range entry.Changes {
			if change.Value.Statuses != nil {
				allStatuses = append(allStatuses, change.Value.Statuses...)
			}
		}
	}
	return allMessages, allStatuses
}
