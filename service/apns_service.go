package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/swimresults/user-service/apns"
	"golang.org/x/net/http2"
)

var apnsUrl = os.Getenv("SR_APNS_URL")

func SendApnsPushNotification(receiver string, title string, subtitle string, message string, interruptionLevel string) (string, string, int, error) {
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
				"interruption-level": "` + interruptionLevel + `",
				"sound": "default"
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
