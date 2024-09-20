package service

import (
	"context"
	"errors"
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var notificationUserCollection *mongo.Collection

func notificationUserService(database *mongo.Database) {
	notificationUserCollection = database.Collection("notification_user")
}

func getNotificationUsersByBsonDocument(d interface{}) ([]model.NotificationUser, error) {
	var users []model.NotificationUser

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := notificationUserCollection.Find(ctx, d)
	if err != nil {
		return []model.NotificationUser{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.NotificationUser
		cursor.Decode(&user)
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return []model.NotificationUser{}, err
	}

	return users, nil
}

func getNotificationUserByBsonDocument(d interface{}) (model.NotificationUser, error) {
	users, err := getNotificationUsersByBsonDocument(d)
	if err != nil {
		return model.NotificationUser{}, err
	}

	if len(users) <= 0 {
		return model.NotificationUser{}, errors.New(entryNotFoundMessage)
	}

	return users[0], nil
}

func GetNotificationUsers() ([]model.NotificationUser, error) {
	return getNotificationUsersByBsonDocument(bson.D{})
}

func GetNotificationUserById(id primitive.ObjectID) (model.NotificationUser, error) {
	return getNotificationUserByBsonDocument(bson.D{{"_id", id}})
}

func RemoveNotificationUserById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := notificationUserCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddNotificationUser(user model.NotificationUser) (model.NotificationUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := notificationUserCollection.InsertOne(ctx, user)
	if err != nil {
		return model.NotificationUser{}, err
	}

	return GetNotificationUserById(r.InsertedID.(primitive.ObjectID))
}

func UpdateNotificationUser(user model.NotificationUser) (model.NotificationUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := notificationUserCollection.ReplaceOne(ctx, bson.D{{"_id", user.Identifier}}, user)
	if err != nil {
		return model.NotificationUser{}, err
	}

	return GetNotificationUserById(user.Identifier)
}

func RegisterNotificationUser(token string) (model.NotificationUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := model.NotificationUser{
		Token: token,
	}

	r, err := notificationUserCollection.InsertOne(ctx, user)
	if err != nil {
		return model.NotificationUser{}, err
	}

	return GetNotificationUserById(r.InsertedID.(primitive.ObjectID))
}
