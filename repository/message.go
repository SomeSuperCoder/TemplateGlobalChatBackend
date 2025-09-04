package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MessageRepo struct {
	Database *mongo.Database
}

func (r *MessageRepo) FindPaged(ctx context.Context, page, limit int64) ([]models.Message, error) {
	var messages []models.Message

	// Set pagination options
	skip := (page - 1) * limit
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(skip)
	opts.SetSort(bson.M{"created_at": -1})

	// Init a cursor
	cursor, err := r.Database.Collection("messages").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepo) CreateMessage(ctx context.Context, message models.Message) error {
	_, err := r.Database.Collection("messages").InsertOne(ctx, message)
	return err
}

func (r *MessageRepo) DeleteMessage(ctx context.Context, messageID bson.ObjectID) error {
	_, err := r.Database.Collection("messages").DeleteOne(ctx, bson.M{"_id": messageID})
	return err
}

func (r *MessageRepo) UpdateMessageText(ctx context.Context, messageID bson.ObjectID, newText string) error {
	update := bson.M{
		"text": newText,
	}

	_, err := r.Database.Collection("messages").UpdateByID(ctx, messageID, update)
	return err
}
