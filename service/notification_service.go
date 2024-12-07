package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/swimresults/user-service/apns"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var apnsUrl = os.Getenv("SR_APNS_URL")

func SendTestPushNotification(receiver string) error {
	token := apns.GetToken()

	t := &http2.Transport{}
	c := &http.Client{
		Transport: t,
	}

	b := []byte(`
		{
			"aps": {
				"alert": {
					"title": "27. IESC 2024",
					"subtitle": "50m Brust männlich",
					"body": "Start in 15 Minuten um ca. 14:34, Lauf 5, Bahn 1"
				}
			}
		}
	`)

	r, err := http.NewRequest("POST", "https://api.sandbox.push.apple.com:443/3/device/"+receiver, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	r.Header.Set("apns-push-type", "alert")
	r.Header.Set("apns-expiration", "0")
	r.Header.Set("apns-priority", "10")
	r.Header.Set("apns-topic", "de.logilutions.SwimResults")

	println("making request with token '" + token + "'")

	resp, err := c.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

	return nil
}

func SendPushNotificationForMeetingAndAthletes(meetingId string, athleteIds []primitive.ObjectID, request dto.MeetingNotificationRequestDto) (int, int, int, error) {
	config, err := GetStoredConfigByMeeting(meetingId)
	if err != nil {
		return 0, 0, 0, err
	}

	if config.Enabled == false {
		fmt.Printf("notifications are diabled for %s\n", meetingId)
		return 0, 0, 0, errors.New("notifications are disabled for " + meetingId)
	}

	meeting, err := GetMeetingById(meetingId)
	if err != nil {
		return 0, 0, 0, err
	}

	athleteIds = append(athleteIds, primitive.NewObjectID())
	athleteIds = append(athleteIds, primitive.NewObjectID())

	var users []model.User
	var err1 error
	if request.MessageType == "athlete" {
		users, err1 = GetUsersByIsMe(athleteIds)
	} else if request.MessageType == "favourites" {
		users, err1 = GetUsersByIsFollower(athleteIds)
	} else {
		users, err1 = GetUsersByIsFollowerOrMe(athleteIds)
	}

	if err1 != nil {
		return 0, 0, 0, err
	}

	var userIds []primitive.ObjectID
	userIds = append(userIds, primitive.NewObjectID())
	userIds = append(userIds, primitive.NewObjectID())
	for _, user := range users {
		userIds = append(userIds, user.Identifier)
	}

	notificationUsers, err := GetNotificationUsersByUserIds(userIds)
	if err != nil {
		return 0, 0, 0, err
	}

	var wg sync.WaitGroup

	success := 0
	for _, user := range notificationUsers {
		if !user.HasSetting(request.MessageType) {
			continue
		}

		wg.Add(1)
		go func(receiver string, title string, subtitle string, message string, interruptionLevel string, success *int) {
			defer wg.Done()
			_, _, code, err := SendPushNotification(receiver, title, subtitle, message, interruptionLevel)
			if err == nil || code == 200 {
				*success++
			}
		}(user.Token, meeting.Series.NameMedium, request.Subtitle, request.Message, request.InterruptionLevel, &success)
	}

	wg.Wait() // Wait for all goroutines to finish

	fmt.Printf("notified %d users with %d/%d devices\n", len(users), success, len(notificationUsers))
	return len(users), len(notificationUsers), success, nil
}

func SendPushNotificationForMeeting(meetingId string, request dto.MeetingNotificationRequestDto) (int, int, int, error) {
	athletes, err := ac.GetAthletesByMeeting(meetingId)
	if err != nil {
		return 0, 0, 0, err
	}

	var athleteIds []primitive.ObjectID
	for _, athlete := range athletes {
		athleteIds = append(athleteIds, athlete.Identifier)
	}

	return SendPushNotificationForMeetingAndAthletes(meetingId, athleteIds, request)
}

func SendPushNotification(receiver string, title string, subtitle string, message string, interruptionLevel string) (string, string, int, error) {
	token := apns.GetToken()

	t := &http2.Transport{}
	c := &http.Client{
		Transport: t,
	}

	if interruptionLevel == "" {
		interruptionLevel = "active"
	}

	b := []byte(`
		{
			"aps": {
				"alert": {
					"title": "` + title + `",
					"subtitle": "` + subtitle + `",
					"body": "` + message + `"
				},
				"interruption-level": "` + interruptionLevel + `"
			}
		}
	`)

	fmt.Printf("notifying user with token: '%s'\n", receiver)

	r, err := http.NewRequest("POST", apnsUrl+"/3/device/"+receiver, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	r.Header.Set("apns-push-type", "alert")
	r.Header.Set("apns-expiration", "0")
	r.Header.Set("apns-priority", "10")
	r.Header.Set("apns-topic", "de.logilutions.SwimResults")

	println("making request with token '" + token + "'")

	resp, err := c.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

	apnsID := resp.Header.Get("apns-unique-id")
	println("apns-unique-id: " + apnsID)

	return apnsID, string(body), resp.StatusCode, nil
}
