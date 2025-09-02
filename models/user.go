package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserSession struct {
	SessionToken string `bson:"session_token"`
	CSRFToken    string `bson:"csrf_token"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Sessions []UserSession      `bson:"sessions"`
}
