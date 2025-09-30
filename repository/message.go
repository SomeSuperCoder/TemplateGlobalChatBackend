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

func (r *MessageRepo) FindPaged(ctx context.Context, page, limit int64) ([]models.Message, int64, error) {
	var messages = []models.Message{}

	// Set pagination options
	skip := (page - 1) * limit
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(skip)
	opts.SetSort(bson.M{"created_at": -1})

	// Init a cursor
	cursor, err := r.Database.Collection("messages").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Extract messages
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, 0, err
	}

	// Get total message count
	count, err := r.Database.Collection("messages").CountDocuments(ctx, bson.M{})

	return messages, count, err // TODO: check if this works)))
}

func (r *MessageRepo) CreateMessage(ctx context.Context, message models.Message) error {
	_, err := r.Database.Collection("messages").InsertOne(ctx, message)
	return err
}

func (r *MessageRepo) DeleteMessage(ctx context.Context, messageID bson.ObjectID) error {
	_, err := r.Database.Collection("messages").DeleteOne(ctx, bson.M{"_id": messageID})
	return err
}

func (r *MessageRepo) UpdateMessage(ctx context.Context, messageID bson.ObjectID, update any) error {
	_, err := r.Database.Collection("messages").UpdateByID(ctx, messageID, bson.M{
		"$set": update,
	})
	return err
}
