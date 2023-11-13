package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Identifier   primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	KeycloakId   string              `json:"keycloak_id,omitempty" bson:"keycloak_id,omitempty"`
	Following    []Following         `json:"following,omitempty" bson:"following,omitempty"`
	OwnAthleteId *primitive.ObjectID `json:"own_athlete_id,omitempty" bson:"own_athlete_id,omitempty"`
	Settings     Settings            `json:"settings,omitempty" bson:"settings,omitempty"`
	Meetings     []string            `json:"meetings,omitempty" bson:"meetings,omitempty"`
}
