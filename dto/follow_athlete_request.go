package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type FollowAthleteRequestDto struct {
	AthleteId primitive.ObjectID `json:"athlete"`
	Follow    bool               `json:"follow,omitempty"`
}
