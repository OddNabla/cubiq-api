package repository

import (
	"context"

	"github.com/DouglasValerio/cubiq-api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InboundMessageRepo struct {
	MongoDatabase *mongo.Database
}

func (r *InboundMessageRepo) InsertInboundMessage(inboundMessage *model.InboundMessage) (*model.InboundMessage, error) {
	collection := r.MongoDatabase.Collection("inbound_messages")
	inboundMessage.SetDefaults()
	_, err := collection.InsertOne(context.Background(), inboundMessage)
	if err != nil {
		return nil, err
	}
	return inboundMessage, nil
}

func (r *InboundMessageRepo) FindInboundMessageById(id string) (*model.InboundMessage, error) {
	collection := r.MongoDatabase.Collection("inbound_messages")
	filter := bson.M{"id": id}
	var inboundMessage model.InboundMessage
	err := collection.FindOne(context.Background(), filter).Decode(&inboundMessage)
	if err != nil {
		return nil, err
	}
	return &inboundMessage, nil
}

func (r *InboundMessageRepo) FindAllInboundMessages() ([]model.InboundMessage, error) {
	collection := r.MongoDatabase.Collection("inbound_messages")
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
