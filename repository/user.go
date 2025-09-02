package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type UserRepo struct {
	Database *mongo.Database
}
