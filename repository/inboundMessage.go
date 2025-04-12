package repository

import (
	"context"
	"time"

	"github.com/DouglasValerio/cubiq-api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InboundMessageRepository interface {
	InsertInboundMessage(inboundMessage *model.InboundMessage) (*model.InboundMessage, error)
	FindInboundMessageById(id string) (*model.InboundMessage, error)
	FindAllInboundMessages() ([]model.InboundMessage, error)
	SetProcessedAt(id string) error
}

type InboundMessageRepo struct {
	MongoDatabase *mongo.Database
}

func (r *InboundMessageRepo) InsertInboundMessage(inboundMessage *model.InboundMessage) (*model.InboundMessage, error) {
	collection := r.MongoDatabase.Collection("InboundMessages")
	inboundMessage.SetDefaults()
	_, err := collection.InsertOne(context.Background(), inboundMessage)
	if err != nil {
		return nil, err
	}
	return inboundMessage, nil
}

func (r *InboundMessageRepo) FindInboundMessageById(id string) (*model.InboundMessage, error) {
	collection := r.MongoDatabase.Collection("InboundMessages")
	filter := bson.M{"id": id}
	var inboundMessage model.InboundMessage
	err := collection.FindOne(context.Background(), filter).Decode(&inboundMessage)
	if err != nil {
		return nil, err
	}
	return &inboundMessage, nil
}

func (r *InboundMessageRepo) FindAllInboundMessages() ([]model.InboundMessage, error) {
	collection := r.MongoDatabase.Collection("InboundMessages")
	result, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var inboundMessages []model.InboundMessage
	if err = result.All(context.Background(), &inboundMessages); err != nil {
		return nil, err
	}
	return inboundMessages, nil
}

func (r *InboundMessageRepo) SetProcessedAt(id string) error {
	collection := r.MongoDatabase.Collection("InboundMessages")
	now := time.Now().UTC()
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"processedAt": now,
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
