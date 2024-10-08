package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type WidgetTile struct {
	WidgetID      primitive.ObjectID `json:"-" bson:"widget_id,omitempty"`
	Widget        Widget             `json:"widget,omitempty" bson:"-"`
	Data          any                `json:"data,omitempty" bson:"data,omitempty"`
	OrderPosition int                `json:"order_position,omitempty" bson:"order_position,omitempty"`
}
