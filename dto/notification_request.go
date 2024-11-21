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
