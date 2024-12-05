package dto

type MeetingNotificationRequestDto struct {
	Subtitle          string `json:"subtitle"`
	Message           string `json:"message"`
	MessageType       string `json:"message_type"`       // like athlete, meeting, schedule or favourites (see settings)
	InterruptionLevel string `json:"interruption_level"` // passive, active, time-sensitive
}

type NotificationRequestDto struct {
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	Message           string `json:"message"`
	InterruptionLevel string `json:"interruption_level"` // passive, active, time-sensitive
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
