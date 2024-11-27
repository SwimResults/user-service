package service

import (
	"bytes"
	"fmt"
	"github.com/swimresults/user-service/apns"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func SendPushBroadcast(channel string, content string) (string, string, string, int, error) {
	token := apns.GetToken()

	t := &http2.Transport{}
	c := &http.Client{
		Transport: t,
	}

	b := []byte(`
		{
			"aps": {
				"timestamp": ` + strconv.Itoa(int(time.Now().Unix())) + `,
				"event": "update",
				"content-state": {
					` + content + `
				}
			}
		}
	`)

	r, err := http.NewRequest("POST", apnsUrl+"/4/broadcasts/apps/de.logilutions.SwimResults", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	r.Header.Set("apns-push-type", "liveactivity")
	r.Header.Set("apns-expiration", "0")
	r.Header.Set("apns-priority", "10")
	r.Header.Set("apns-channel-id", channel)

	println("making broadcast request with token '" + token + "'")

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

	apnsRequestID := resp.Header.Get("apns-request-id")
	apnsUID := resp.Header.Get("apns-unique-id")
	println("apns-request-id: " + apnsRequestID)
	println("apns-unique-id: " + apnsUID)

	return apnsRequestID, apnsUID, string(body), resp.StatusCode, nil
}