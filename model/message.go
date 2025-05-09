package model

import "time"

type ChatMessage struct {
	Id        string                  `json:"id" bson:"_id"`
	From      string                  `json:"from" bson:"from"`
	Timestamp time.Time               `json:"timestamp" bson:"timestamp"`
	Type      string                  `json:"type" bson:"type"`
	Text      MessageText             `json:"text" bson:"text"`
	Media     MediaMessage            `json:"media" bson:"media"`
	Location  LocationMessage         `json:"location" bson:"location"`
	Contacts  []SharedContactsMessage `json:"contacts" bson:"contacts"`
	Reaction  ReactionMessage         `json:"reaction" bson:"reaction"`
	Statuses  []MessageStatus         `json:"statuses" bson:"statuses"`
	CreatedAt time.Time               `json:"createdAt" bson:"createdAt"`
	Summary   string                  `json:"summary" bson:"summary"`
}

type MessageStatus struct {
	Id        string `json:"id" bson:"id"`
	Status    string `json:"status" bson:"status"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
}

type MediaMessage struct {
	Caption  string `json:"caption" bson:"caption"`
	Filename string `json:"filename" bson:"filename"`
	Id       string `json:"id" binding:"required" bson:"id"`
	MimeType string `json:"mime_type" binding:"required" bson:"mime_type"`
	Sha256   string `json:"sha256" binding:"required" bson:"sha256"`
	FileUrl  string `json:"file_url" bson:"file_url"`
}

func (c *ChatMessage) SetMediaFileUrl(fileUrl string) {
	c.Media.FileUrl = fileUrl
}

type SharedContactsMessage struct {
	Name   SharedContactNameMessage    `json:"name" binding:"required" bson:"name"`
	Phones []SharedContactPhoneMessage `json:"phones" binding:"required" bson:"phones"`
}

type SharedContactNameMessage struct {
	FirstName     string `json:"first_name" binding:"required" bson:"first_name"`
	FormattedName string `json:"formatted_name" binding:"required" bson:"formatted_name"`
	LastName      string `json:"last_name" binding:"required" bson:"last_name"`
	MiddleName    string `json:"middle_name" binding:"required" bson:"middle_name"`
}

type SharedContactPhoneMessage struct {
	Phone string `json:"phone" binding:"required" bson:"phone"`
	Type  string `json:"type" binding:"required" bson:"type"`
	WaId  string `json:"wa_id" binding:"required" bson:"wa_id"`
}

type ReactionMessage struct {
	MessageId string `json:"message_id" binding:"required" bson:"message_id"`
	Emoji     string `json:"emoji" binding:"required" bson:"emoji"`
}
type LocationMessage struct {
	Address   string  `json:"address" bson:"address"`
	Latitude  float32 `json:"latitude" bson:"latitude"`
	Longitude float32 `json:"longitude" bson:"longitude"`
	Name      string  `json:"name" bson:"name"`
}

func (m *ChatMessage) SetDefaults() {
	now := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
}
func (m *ChatMessage) SetSummary() {
	if m.Summary != "" {
		return
	}
	if m.Type == "text" {
		m.Summary = m.Text.Body
		return
	}
	switch m.Type {
	case "media":
		m.Summary = m.Type

	case "location":
		m.Summary = "Localização"

	case "contacts":
		if len(m.Contacts) > 0 {
			m.Summary = "Contato: " + m.Contacts[0].Name.FormattedName
		} else {
			m.Summary = "Contatos"
		}

	case "reaction":
		m.Summary = "Reação: " + m.Reaction.Emoji
	case "image":
		m.Summary = "Imagem"
	case "video":
		m.Summary = "Vídeo"
	case "audio":
		m.Summary = "Áudio"
	case "document":
		m.Summary = "Documento"
	case "sticker":
		m.Summary = "Sticker"
	default:
		m.Summary = "Nova mensagem"
	}
}
