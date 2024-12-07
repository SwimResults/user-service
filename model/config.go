package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Config struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Meeting    string             `json:"meeting,omitempty" bson:"meeting,omitempty"`
	Enabled    bool               `json:"enabled,omitempty" bson:"enabled,omitempty"`
}
