package dto

type NotificationRequestDto struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Message  string `json:"message"`
}

type NotificationResponseDto struct {
	ApnsId string `json:"apns_id"`
	Body   string `json:"body"`
}

type NotificationsResponseDto struct {
	UserCount             int `json:"user_count"`
	NotificationUserCount int `json:"notification_user_count"`
	SuccessCount          int `json:"success_count"`
}
