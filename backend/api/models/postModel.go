package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostModel struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Creator      primitive.ObjectID `json:"creator" bson:"creator"`
	Title        string             `json:"title" bson:"title"`
	Message      string             `json:"message" bson:"message"`
	SelectedFile string             `json:"selectedFile" bson:"selectedFile"`
	Likes        []string           `json:"likes" bson:"likes"`
	Comments     []CommentWithUser  `json:"comments,omitempty" bson:"comments,omitempty"`
	User         *UserModel         `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

type CommentWithUser struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"postId,omitempty" bson:"postId,omitempty"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Value     string             `json:"value" bson:"value"`
	CreatedAt time.Time          `json:"craetedAt" bson:"createdAt"`
	User      UserModel          `json:"user" bson:"user"`
}

// interfaces
type CreateOrUpdatePost struct {
	Title        string `json:"title" bson:"title" validate:"required"`
	Message      string `json:"message" bson:"message" validate:"required,min=5"`
	SelectedFile string `json:"selectedFile" bson:"selectedFile"`
}

// interfaces
type ComnmentPost struct {
	Value string `json:"value" bson:"value" validate:"required"`
}
