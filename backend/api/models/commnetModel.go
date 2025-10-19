package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"postId,omitempty" bson:"postId,omitempty"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Value     string             `json:"value" bson:"value"`
	CreatedAt time.Time          `json:"craetedAt" bson:"createdAt"`
}

// interfaces
type CreateComment struct {
	Value string `json:"value" bson:"value" validate:"required"`
}
