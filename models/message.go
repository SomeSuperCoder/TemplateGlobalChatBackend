package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Author   bson.ObjectID `bson:"author" json:"author"`
	Text     string        `bson:"text" json:"text"`
	CratedAt time.Time     `bson:"created_at" json:"created_at"`
}
