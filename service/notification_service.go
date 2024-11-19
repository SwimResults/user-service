package service

import (
	"bytes"
	"fmt"
	"github.com/swimresults/user-service/apns"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
)

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
