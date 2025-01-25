package service

import (
	"context"
	"errors"
	"fmt"
	athleteClient "github.com/swimresults/athlete-service/client"
	meetingClient "github.com/swimresults/meeting-service/client"
	meetingModel "github.com/swimresults/meeting-service/model"
	"github.com/swimresults/user-service/model"
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

var meetings = make(map[string]meetingModel.Meeting)
var configs = make(map[string]model.Config)

func Init(c *mongo.Client) {
	database := c.Database(os.Getenv("SR_USER_MONGO_DATABASE"))
	client = c

	InitMeetings()

	userService(database)
	widgetService(database)
	dashboardService(database)
	notificationUserService(database)
	configService(database)

	InitConfigs()
}

func GetMeetingById(id string) (*meetingModel.Meeting, error) {
	existing := meetings[id]
	if existing.MeetId == id {
		fmt.Printf("returning cached meeting: %s\nvalues:\nid: %s;\nstart_date: %s;\nftp_mask: %s;\npush_channel: %s;\n", id, existing.MeetId, existing.DateStart, existing.Data.FtpStartListMask, existing.Data.PushNotificationChannel)
		return &existing, nil
	}
	meeting, err := mc.GetMeetingById(id)
	if err != nil {
		return nil, err
	}

	meetings[id] = *meeting

	return meeting, nil
}

func GetStoredConfigByMeeting(meeting string) (*model.Config, error) {
	existing := configs[meeting]
	if existing.Meeting == meeting {
		fmt.Printf("returning cached config: %s\n", meeting)
		return &existing, nil
	}

	return nil, errors.New("no config found")
}

func UpdateStoredConfigByMeeting(meeting string, config model.Config) {
	configs[meeting] = config
}

func InitMeetings() {
	meetingList, err := mc.GetMeetings()
	if err != nil {
		fmt.Printf("Failed loading meetings: %s", err.Error())
		return
	}

	for _, meeting := range meetingList {
		fmt.Printf("meeting: %s, channel: %s\n", meeting.MeetId, meeting.Data.PushNotificationChannel)
		meetings[meeting.MeetId] = meeting
	}

	println("added:")
	for key, meeting := range meetings {
		fmt.Printf("%s: %s\n", key, meeting.MeetId)
	}
}

func InitConfigs() {
	configList, err := GetConfigs()
	if err != nil {
		fmt.Printf("Failed loading meetings: %s", err.Error())
		return
	}

	for _, config := range configList {
		configs[config.Meeting] = config
	}
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
