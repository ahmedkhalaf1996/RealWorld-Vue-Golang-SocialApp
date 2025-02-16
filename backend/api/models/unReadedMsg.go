package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnReadedMsg struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MainUserid            string             `json:"mainUserid" bson:"mainUserid"`
	OtherUserid           string             `json:"otherUserid" bson:"otherUserid"`
	NumOfUnreadedMessages int                `json:"numOfUnreadedMessages" bson:"numOfUnreadedMessages"`
	IsReaded              bool               `json:"isReaded" bson:"isReaded"`
}
