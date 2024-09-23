package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FollowingAthlete struct {
	AthleteId primitive.ObjectID `json:"athlete_id" bson:"athlete_id"`
	AddedAt   time.Time          `json:"added_at,omitempty" bson:"added_at,omitempty"`
}

// FollowingTeam TODO make teams followable
type FollowingTeam struct {
	TeamId  primitive.ObjectID `json:"team_id" bson:"team_id"`
	AddedAt time.Time          `json:"added_at,omitempty" bson:"added_at,omitempty"`
}

// Following TODO use merged struct
type Following struct {
	FollowingAthletes []FollowingAthlete `json:"following_athletes,omitempty" bson:"following_athletes,omitempty"`
	FollowingTeams    []FollowingTeam    `json:"following_teams,omitempty" bson:"following_teams,omitempty"`
}
