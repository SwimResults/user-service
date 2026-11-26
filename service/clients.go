package service

// Clients for meeting- and athlete-service, built on the generated packages from
// github.com/swimresults/swimresults-api. They keep the method set and semantics of the per-service
// client packages they replace.

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/swimresults/swimresults-api/go/athleteapi"
	"github.com/swimresults/swimresults-api/go/meetingapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const requestTimeout = 30 * time.Second

var httpClient = &http.Client{Timeout: requestTimeout}

// serviceKey authenticates user-service against the other services, like the service-core client did.
func serviceKey(_ context.Context, req *http.Request) error {
	if key := os.Getenv("SR_SERVICE_KEY"); key != "" {
		req.Header.Set("X-Swimresults-Service", key)
	}
	return nil
}

func reqCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

// ---------------------------------------------------------------- meeting-service

type MeetingClient struct {
	c *meetingapi.ClientWithResponses
}

func newMeetingClient(url string) *MeetingClient {
	c, err := meetingapi.NewClientWithResponses(strings.TrimRight(url, "/"),
		meetingapi.WithHTTPClient(httpClient), meetingapi.WithRequestEditorFn(serviceKey))
	if err != nil {
		fmt.Printf("invalid meeting-service URL %q: %v\n", url, err)
	}
	return &MeetingClient{c}
}

func (m *MeetingClient) GetMeetings() ([]meetingapi.Meeting, error) {
	if m.c == nil {
		return nil, fmt.Errorf("meeting client is not configured")
	}
	ctx, cancel := reqCtx()
	defer cancel()
	res, err := m.c.GetMeetingsWithResponse(ctx, &meetingapi.GetMeetingsParams{})
	if err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("GetMeetings received error: %d", res.StatusCode())
	}
	return *res.JSON200, nil
}

func (m *MeetingClient) GetMeetingById(id string) (*meetingapi.Meeting, error) {
	if m.c == nil {
		return nil, fmt.Errorf("meeting client is not configured")
	}
	ctx, cancel := reqCtx()
	defer cancel()
	res, err := m.c.GetMeetingByMeetIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("GetMeetingById received error: %d", res.StatusCode())
	}
	return res.JSON200, nil
}

// ---------------------------------------------------------------- athlete-service

type AthleteClient struct {
	c *athleteapi.ClientWithResponses
}

func newAthleteClient(url string) *AthleteClient {
	c, err := athleteapi.NewClientWithResponses(strings.TrimRight(url, "/"),
		athleteapi.WithHTTPClient(httpClient), athleteapi.WithRequestEditorFn(serviceKey))
	if err != nil {
		fmt.Printf("invalid athlete-service URL %q: %v\n", url, err)
	}
	return &AthleteClient{c}
}

// ForeignAthlete is the part of an athlete-service athlete that user-service uses.
type ForeignAthlete struct {
	Identifier primitive.ObjectID
}

// GetAthletesByMeeting returns all athletes of a meeting (no paging), nil if the meeting has none (404).
func (a *AthleteClient) GetAthletesByMeeting(meeting string) ([]ForeignAthlete, error) {
	if a.c == nil {
		return nil, fmt.Errorf("athlete client is not configured")
	}
	ctx, cancel := reqCtx()
	defer cancel()
	res, err := a.c.GetAthletesByMeetingWithResponse(ctx, meeting, &athleteapi.GetAthletesByMeetingParams{})
	if err != nil {
		return nil, err
	}
	switch {
	case res.StatusCode() == http.StatusNotFound:
		return nil, nil
	case res.JSON200 == nil:
		return nil, fmt.Errorf("GetAthletesByMeeting received error: %d", res.StatusCode())
	}
	athletes := make([]ForeignAthlete, 0, len(*res.JSON200))
	for _, at := range *res.JSON200 {
		id, err := primitive.ObjectIDFromHex(at.Identifier)
		if err != nil && at.Identifier != "" {
			return nil, fmt.Errorf("GetAthletesByMeeting: invalid athlete id %q: %w", at.Identifier, err)
		}
		athletes = append(athletes, ForeignAthlete{Identifier: id})
	}
	return athletes, nil
}
