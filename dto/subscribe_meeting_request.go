package dto

type SubscribeMeetingRequestDto struct {
	Meeting   string `json:"meeting"`
	Subscribe bool   `json:"subscribe,omitempty"`
}
