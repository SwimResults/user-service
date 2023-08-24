package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Following struct {
	AthleteId primitive.ObjectID `json:"athlete_id" bson:"athlete_id"`
	AddedAt   time.Time          `json:"added_at,omitempty" bson:"added_at,omitempty"`
}
