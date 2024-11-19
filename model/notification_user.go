package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationUser struct {
	Identifier primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId     primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Token      string             `json:"token" bson:"token"`
	Settings   Settings           `json:"settings" bson:"settings"`
	Device     Device             `json:"device" bson:"device"`
	Meetings   []string           `json:"meetings" bson:"meetings"`
}
