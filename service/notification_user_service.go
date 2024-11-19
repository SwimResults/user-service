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

func GetNotificationUserByToken(token string) (model.NotificationUser, error) {
	return getNotificationUserByBsonDocument(bson.D{{"token", token}})
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

// RegisterNotificationUser adds the given token to the database,
// if the token already exists, it just returns the user
func RegisterNotificationUser(token string, device model.Device, user *model.User) (model.NotificationUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, err := GetNotificationUserByToken(token)
	if err != nil {
		if err.Error() == entryNotFoundMessage {
			notificationUser := model.NotificationUser{
				Token:  token,
				Device: device,
			}

			if user != nil {
				notificationUser.UserId = user.Identifier
			}

			r, err2 := notificationUserCollection.InsertOne(ctx, notificationUser)
			if err2 != nil {
				return model.NotificationUser{}, err2
			}

			existing, _ = GetNotificationUserById(r.InsertedID.(primitive.ObjectID))
		} else {
			return model.NotificationUser{}, err
		}
	} else {
		existing.Device = device
		if user != nil {
			existing.UserId = user.Identifier
		}

		r, err2 := notificationUserCollection.ReplaceOne(ctx, bson.D{{"_id", existing.Identifier}}, existing)
		if err2 != nil {
			return model.NotificationUser{}, err2
		}

		existing, _ = GetNotificationUserById(r.UpsertedID.(primitive.ObjectID))
	}

	return existing, nil
}
