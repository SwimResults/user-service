package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type req struct{ method, uri, serviceKey string }

func stub(t *testing.T, status int, body string) (string, *req) {
	t.Helper()
	r := &req{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, hr *http.Request) {
		r.method, r.uri, r.serviceKey = hr.Method, hr.URL.RequestURI(), hr.Header.Get("X-Swimresults-Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv.URL + "/", r // env URLs end with "/"
}

const meetingJSON = `{"meet_id":"IESC19","date_start":"2019-06-01T00:00:00Z",
  "series":{"name_medium":"Erzgebirgs-Schwimmcup"},
  "data":{"push_notification_channel":"chan-1","ftp_start_list_mask":"mask"}}`

func TestGetMeetings(t *testing.T) {
	t.Setenv("SR_SERVICE_KEY", "k")
	url, r := stub(t, 200, "["+meetingJSON+"]")
	meetings, err := newMeetingClient(url).GetMeetings()
	if err != nil || len(meetings) != 1 {
		t.Fatalf("meetings=%v err=%v", meetings, err)
	}
	m := meetings[0]
	if m.MeetId != "IESC19" || m.Data.PushNotificationChannel != "chan-1" || m.Data.FtpStartListMask != "mask" || m.Series.NameMedium != "Erzgebirgs-Schwimmcup" {
		t.Errorf("meeting = %+v", m)
	}
	if r.method != http.MethodGet || r.uri != "/meeting" || r.serviceKey != "k" {
		t.Errorf("request = %s %s key=%q", r.method, r.uri, r.serviceKey)
	}
}

func TestGetMeetingById(t *testing.T) {
	url, r := stub(t, 200, meetingJSON)
	m, err := newMeetingClient(url).GetMeetingById("IESC19")
	if err != nil || m.MeetId != "IESC19" || r.uri != "/meeting/meet/IESC19" {
		t.Fatalf("m=%+v err=%v uri=%s", m, err, r.uri)
	}
	url, _ = stub(t, 404, `{"message":"not found"}`)
	if _, err := newMeetingClient(url).GetMeetingById("X"); err == nil {
		t.Error("404: want error, like before")
	}
}

func TestGetAthletesByMeeting(t *testing.T) {
	url, r := stub(t, 200, `[{"_id":"650c2e653df25edacb51bd3b","name":"A"},{"_id":"650c2e653df25edacb51bd3c","name":"B"}]`)
	athletes, err := newAthleteClient(url).GetAthletesByMeeting("IESC19")
	if err != nil || len(athletes) != 2 || athletes[1].Identifier.Hex() != "650c2e653df25edacb51bd3c" {
		t.Fatalf("athletes=%v err=%v", athletes, err)
	}
	if r.uri != "/athlete/meet/IESC19" { // no paging params: the service returns all athletes
		t.Errorf("uri = %s", r.uri)
	}

	url, _ = stub(t, 404, `{"message":"not found"}`)
	if athletes, err := newAthleteClient(url).GetAthletesByMeeting("X"); err != nil || athletes != nil {
		t.Fatalf("404: athletes=%v err=%v, want nil,nil", athletes, err)
	}

	url, _ = stub(t, 200, `[{"_id":"bad"}]`)
	if _, err := newAthleteClient(url).GetAthletesByMeeting("X"); err == nil {
		t.Error("invalid id: want error (the old client failed decoding)")
	}
}
