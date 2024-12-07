package service

import (
	"context"
	"errors"
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var configCollection *mongo.Collection

func configService(database *mongo.Database) {
	configCollection = database.Collection("config")
}

func getConfigsByBsonDocument(d interface{}) ([]model.Config, error) {
	var configs []model.Config

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := configCollection.Find(ctx, d)
	if err != nil {
		return []model.Config{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var config model.Config
		cursor.Decode(&config)
		configs = append(configs, config)
	}

	if err := cursor.Err(); err != nil {
		return []model.Config{}, err
	}

	return configs, nil
}

func getConfigByBsonDocument(d interface{}) (model.Config, error) {
	configs, err := getConfigsByBsonDocument(d)
	if err != nil {
		return model.Config{}, err
	}

	if len(configs) <= 0 {
		return model.Config{}, errors.New(entryNotFoundMessage)
	}

	return configs[0], nil
}

func GetConfigByMeeting(meeting string) (model.Config, error) {
	return getConfigByBsonDocument(bson.M{"meeting": meeting})
}

func GetConfigs() ([]model.Config, error) {
	return getConfigsByBsonDocument(bson.M{})
}

func AddConfig(config model.Config) (model.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := configCollection.InsertOne(ctx, config)
	if err != nil {
		return model.Config{}, err
	}

	return GetConfigByMeeting(config.Meeting)
}

func DisableNotification(meeting string, enabled bool) (*model.Config, error) {
	config, err := GetStoredConfigByMeeting(meeting)
	if err != nil {
		return nil, err
	}

	config.Enabled = enabled

	UpdateStoredConfigByMeeting(meeting, *config)
	newConfig, err := UpdateConfig(*config)
	if err != nil {
		return nil, err
	}

	return &newConfig, nil
}

func UpdateConfig(config model.Config) (model.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.D{{"meeting", config.Meeting}}, config)
	if err != nil {
		return model.Config{}, err
	}

	return GetConfigByMeeting(config.Meeting)
}
