package dto

import "github.com/swimresults/user-service/model"

type RegisterNotificationUserRequestDto struct {
	Token  string       `json:"token,omitempty"`
	Device model.Device `json:"device,omitempty"`
}
