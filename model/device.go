package model

type Device struct {
	Name           string `json:"name" bson:"name"`
	Model          string `json:"model" bson:"model"`
	LocalizedModel string `json:"localized_model" bson:"localized_model"`
	SystemName     string `json:"system_name" bson:"system_name"`
	SystemVersion  string `json:"system_version" bson:"system_version"`
	Type           string `json:"type" bson:"type"`
	UISize         string `json:"ui_size" bson:"ui_size"`
	AppVersion     string `json:"app_version" bson:"app_version"`
}
