package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationUser struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Token      string             `json:"token,omitempty" bson:"token,omitempty"`
	Settings   Settings           `json:"settings,omitempty" bson:"settings,omitempty"`
	Meetings   []string           `json:"meetings,omitempty" bson:"meetings,omitempty"`
}
