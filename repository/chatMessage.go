package repository

import (
	"context"

	"github.com/DouglasValerio/cubiq-api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatMessageRepository interface {
	InsertChatMessage(chatMessage *model.ChatMessage) (*model.ChatMessage, error)
	UpdateChatMessageStatus(id string, statuses []model.MessageStatus) error
}

type ChatMessageRepo struct {
	MongoDatabase *mongo.Database
}

func (r *ChatMessageRepo) InsertChatMessage(chatMessage *model.ChatMessage) (*model.ChatMessage, error) {
	collection := r.MongoDatabase.Collection("chat_messages")
	_, err := collection.InsertOne(context.Background(), chatMessage)
	if err != nil {
		return nil, err
	}
	return chatMessage, nil
}

func (r *ChatMessageRepo) UpdateChatMessageStatus(id string, statuses []model.MessageStatus) error {
	collection := r.MongoDatabase.Collection("chat_messages")
	filter := bson.M{"_id": id}
	update := bson.M{
		"$push": bson.M{
			"statuses": bson.M{
				"$each": statuses,
			},
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
