package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swimresults/service-core/client"
	"github.com/swimresults/user-service/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationClient struct {
	apiUrl string
}

func NewNotificationClient(url string) *NotificationClient {
	return &NotificationClient{apiUrl: url}
}

func (c *NotificationClient) SendNotification(title string, subtitle string, message string) (*dto.NotificationResponseDto, error) {
	request := dto.NotificationRequestDto{
		Title:    title,
		Subtitle: subtitle,
		Message:  message,
	}

	res, err := client.Post(c.apiUrl, "notification/import", request, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseDto := &dto.NotificationResponseDto{}
	err = json.NewDecoder(res.Body).Decode(responseDto)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notification request returned: %d", res.StatusCode)
	}
	return responseDto, nil
}

func (c *NotificationClient) SendNotificationForMeetingAndAthlete(meeting string, athleteId primitive.ObjectID, subtitle string, message string, messageType string, interruptionLevel string) (*dto.NotificationResponseDto, error) {
	request := dto.MeetingNotificationRequestDto{
		Subtitle:          subtitle,
		Message:           message,
		MessageType:       messageType,
		InterruptionLevel: interruptionLevel,
	}

	res, err := client.Post(c.apiUrl, "notification/meet/"+meeting+"/athlete/"+athleteId.Hex(), request, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseDto := &dto.NotificationResponseDto{}
	err = json.NewDecoder(res.Body).Decode(responseDto)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notification request for meeting and athlete returned: %d", res.StatusCode)
	}
	return responseDto, nil
}

func (c *NotificationClient) SendMeetingBroadcastNotification(meeting string, body interface{}) (*dto.BroadcastResponseDto, error) {
	fmt.Printf("sending meeting broadcast request to: '%s'\n", "notification/broadcast/meeting/"+meeting)

	res, err := client.Post(c.apiUrl, "notification/broadcast/meeting/"+meeting, body, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseDto := &dto.BroadcastResponseDto{}
	err = json.NewDecoder(res.Body).Decode(responseDto)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meeting broadcast request returned: %d", res.StatusCode)
	}
	return responseDto, nil
}
