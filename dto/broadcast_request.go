package dto

type BroadcastResponseDto struct {
	ApnsRequestId string `json:"apns_request_id"`
	ApnsUniqueId  string `json:"apns_unique_id"`
	Body          string `json:"body"`
}
