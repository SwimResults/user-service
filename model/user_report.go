package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReportSubjectType string

const (
	ReportTypeStart   ReportSubjectType = "START"
	ReportTypeHeat    ReportSubjectType = "HEAT"
	ReportTypeEvent   ReportSubjectType = "EVENT"
	ReportTypeMeeting ReportSubjectType = "MEETING"
	ReportTypeAthlete ReportSubjectType = "ATHLETE"
	ReportTypeTeam    ReportSubjectType = "TEAM"
)

type UserReport struct {
	Identifier   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Message      string             `json:"message" bson:"message"`
	Meeting      string             `json:"meeting,omitempty" bson:"meeting,omitempty"`
	UserId       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	SubjectId    primitive.ObjectID `json:"subject_id,omitempty" bson:"subject_id,omitempty"`
	SubjectType  ReportSubjectType  `json:"subject_type,omitempty" bson:"subject_type,omitempty"`
	Acknowledged bool               `json:"acknowledged" bson:"acknowledged"`
	Completed    bool               `json:"completed" bson:"completed"`
	AddedAt      time.Time          `json:"added_at,omitempty" bson:"added_at,omitempty"`
}

func GetReportSubjectTypes() []ReportSubjectType {
	return []ReportSubjectType{
		ReportTypeStart,
		ReportTypeHeat,
		ReportTypeEvent,
		ReportTypeMeeting,
		ReportTypeAthlete,
		ReportTypeTeam,
	}
}
