package service

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var app *firebase.App

func fcmService() {
	opt := option.WithCredentialsFile("config/fcm/private_key.json")

	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		fmt.Println("Failed to initialized fcm app")
	}

}

func SendFcmPushNotification(receiver string, title string, subtitle string, body string, interruptionLevel string) (string, string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Replace with actual implementation and error handling
	// Initialize Firebase app using credentials
	fmcClient, _ := app.Messaging(ctx)

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  subtitle + ": " + body,
		},
		Data: map[string]string{
			"interruptionLevel": interruptionLevel,
			"subtitle":          subtitle,
		},
		Token: receiver, // The FCM token retrieved from the Android device
	}

	response, err := fmcClient.Send(ctx, message)

	return response, "", 0, err
}
