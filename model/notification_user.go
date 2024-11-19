package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationUser struct {
	Identifier primitive.ObjectID `json:"_id" bson:"_id"`
	UserId     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Token      string             `json:"token" bson:"token"`
	Settings   Settings           `json:"settings" bson:"settings"`
	Device     Device             `json:"device" bson:"device"`
	Meetings   []string           `json:"meetings" bson:"meetings"`
}
