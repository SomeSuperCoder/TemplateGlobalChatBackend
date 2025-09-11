package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserSession struct {
	SessionToken string    `bson:"session_token" json:"session_token"`
	CSRFToken    string    `bson:"csrf_token" json:"csrf_token"`
	CratedAt     time.Time `bson:"created_at" json:"created_at"`
}

type User struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Username       string        `bson:"username" json:"username"`
	HashedPassword string        `bson:"hashed_password" json:"hashed_password"`
	Sessions       []UserSession `bson:"sessions" json:"sessions"`
	CratedAt       time.Time     `bson:"created_at" json:"created_at"`
}
