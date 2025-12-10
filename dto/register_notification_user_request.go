package dto

import "github.com/swimresults/user-service/model"

type RegisterNotificationUserRequestDto struct {
	Token       string                     `json:"token,omitempty"`
	Device      model.Device               `json:"device,omitempty"`
	Settings    model.NotificationSettings `json:"settings,omitempty"`
	PushService string                     `json:"push_service" bson:"push_service"` // PushService: APNS, FCM
}
