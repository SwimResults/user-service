package model

type Settings struct {
	Language string `json:"language,omitempty" bson:"language,omitempty"`
	Theme    string `json:"theme,omitempty" bson:"theme,omitempty"`
}
