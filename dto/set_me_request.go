package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type SetMeRequestDto struct {
	AthleteId primitive.ObjectID `json:"athlete"`
	Set       bool               `json:"set,omitempty"`
}
