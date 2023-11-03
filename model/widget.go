package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Widget struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Size       string             `json:"size,omitempty" bson:"size,omitempty"`
	Content    string             `json:"content,omitempty" bson:"content,omitempty"`
}
