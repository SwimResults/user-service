package model

type WidgetContainer struct {
	OrderPosition int          `json:"order_position,omitempty" bson:"order_position,omitempty"`
	Widgets       []WidgetTile `json:"widgets,omitempty" bson:"widgets,omitempty"`
}
