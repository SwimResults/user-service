package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NotificationUser struct {
	Identifier primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	UserId     primitive.ObjectID   `json:"user_id" bson:"user_id,omitempty"`
	Token      string               `json:"token" bson:"token"`
	Settings   NotificationSettings `json:"settings" bson:"settings"`
	Device     Device               `json:"device" bson:"device"`
	AddedAt    time.Time            `json:"added_at,omitempty" bson:"added_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type NotificationSettings struct {
	Athlete    bool `json:"token" bson:"token"`
	Favourites bool `json:"favourites" bson:"favourites"`
	Meeting    bool `json:"meeting" bson:"meeting"`
	Schedule   bool `json:"schedule" bson:"schedule"`
}
