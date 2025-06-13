package dto

import (
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserReportSubmission struct {
	Message     string                  `json:"message"`
	Anonymous   bool                    `json:"anonymous,omitempty"`
	Meeting     string                  `json:"meeting,omitempty"`
	SubjectId   primitive.ObjectID      `json:"subject_id,omitempty"`
	SubjectType model.ReportSubjectType `json:"subject_type,omitempty"`
}
