package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	author   bson.ObjectID `bson:"author"`
	text     string        `bson:"text"`
	CratedAt time.Time     `bson:"created_at"`
}
