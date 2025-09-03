package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserSession struct {
	SessionToken string    `bson:"session_token"`
	CSRFToken    string    `bson:"csrf_token"`
	CratedAt     time.Time `bson:"created_at"`
}

type User struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	Username       string        `bson:"username"`
	HashedPassword string        `bson:"hashed_password"`
	Sessions       []UserSession `bson:"sessions"`
	CratedAt       time.Time     `bson:"created_at"`
}
