package service

import (
	"context"
	"fmt"
	athleteClient "github.com/swimresults/athlete-service/client"
	meetingClient "github.com/swimresults/meeting-service/client"
	meetingModel "github.com/swimresults/meeting-service/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var meetingServiceUrl = os.Getenv("SR_USER_MEETING_URL")
var athleteServiceUrl = os.Getenv("SR_USER_ATHLETE_URL")

var mc = meetingClient.NewMeetingClient(meetingServiceUrl)
var ac = athleteClient.NewAthleteClient(athleteServiceUrl)

var client *mongo.Client

var meetings = make(map[string]*meetingModel.Meeting)

func Init(c *mongo.Client) {
	database := c.Database(os.Getenv("SR_USER_MONGO_DATABASE"))
	client = c

	userService(database)
	widgetService(database)
	dashboardService(database)
	notificationUserService(database)
}

func GetMeetingById(id string) (*meetingModel.Meeting, error) {
	existing := meetings[id]
	if existing != nil {
		fmt.Printf("returning cached meeting: %s", id)
		return existing, nil
	}
	meeting, err := mc.GetMeetingById(id)
	if err != nil {
		return nil, err
	}

	meetings[id] = meeting

	return meeting, nil
}

func PingDatabase() bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	return true
}
