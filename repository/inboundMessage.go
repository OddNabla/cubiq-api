package repository

import (
	"context"

	"github.com/DouglasValerio/cubiq-api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InboundMessageRepo struct {
	MongoCollection *mongo.Collection
}

func (r *InboundMessageRepo) InsertInboundMessage(inboundMessage *model.InboundMessage) (*model.InboundMessage, error) {
	inboundMessage.SetDefaults()
	_, err := r.MongoCollection.InsertOne(context.Background(), inboundMessage)
	if err != nil {
		return nil, err
	}
	return inboundMessage, nil
}

func (r *InboundMessageRepo) FindInboundMessageById(id string) (*model.InboundMessage, error) {
	filter := bson.M{"id": id}
	var inboundMessage model.InboundMessage
	err := r.MongoCollection.FindOne(context.Background(), filter).Decode(&inboundMessage)
	if err != nil {
		return nil, err
	}
	return &inboundMessage, nil
}

func (r *InboundMessageRepo) FindAllInboundMessages() ([]model.InboundMessage, error) {
	result, err := r.MongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var inboundMessages []model.InboundMessage
	if err = result.All(context.Background(), &inboundMessages); err != nil {
		return nil, err
	}
	return inboundMessages, nil
}
