package repository

import (
	"context"
	"errors"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	Database *mongo.Database
}

func (r *UserRepo) Register(ctx context.Context, user *models.User) {
	r.Database.Collection("users").InsertOne(ctx, user)
}

func (r *UserRepo) DoesExist(ctx context.Context, username string) bool {
	res := r.Database.Collection("users").FindOne(ctx, bson.D{{Key: "username", Value: username}})
	return !errors.Is(res.Err(), mongo.ErrNoDocuments)
}
