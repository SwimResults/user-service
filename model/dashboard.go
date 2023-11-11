package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dashboard struct {
	Identifier      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	WidgetContainer []WidgetContainer  `json:"widget_container,omitempty" bson:"widget_container,omitempty"`
	User            string             `json:"user,omitempty" bson:"user,omitempty"`
	Official        bool               `json:"official,omitempty" bson:"official,omitempty"`
	Default         bool               `json:"default,omitempty" bson:"default,omitempty"`
	MeetingStates   []string           `json:"meeting_states,omitempty" bson:"meeting_states,omitempty"`
}
