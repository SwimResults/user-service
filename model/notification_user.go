package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationUser struct {
	Identifier  primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	UserId      primitive.ObjectID   `json:"user_id" bson:"user_id,omitempty"`
	Token       string               `json:"token" bson:"token"`
	Settings    NotificationSettings `json:"settings,omitempty" bson:"settings"`
	Device      Device               `json:"device" bson:"device"`
	PushService string               `json:"push_service" bson:"push_service"` // PushService: APNS, FCM
	AddedAt     time.Time            `json:"added_at,omitempty" bson:"added_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type NotificationSettings struct {
	Athlete    bool `json:"athlete" bson:"athlete"`
	Favourites bool `json:"favourites" bson:"favourites"`
	Meeting    bool `json:"meeting" bson:"meeting"`
	Schedule   bool `json:"schedule" bson:"schedule"`
}

func (user *NotificationUser) HasSetting(notificationType string) bool {
	switch notificationType {
	case "athlete":
		return user.Settings.Athlete
	case "favourites":
		return user.Settings.Favourites
	case "meeting":
		return user.Settings.Meeting
	case "schedule":
		return user.Settings.Schedule

	}
	return false
}
