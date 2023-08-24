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

var collection *mongo.Collection

func userService(database *mongo.Database) {
	collection = database.Collection("user")
}

func getUsersByBsonDocument(d interface{}) ([]model.User, error) {
	var users []model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, d)
	if err != nil {
		return []model.User{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func getUserByBsonDocument(d interface{}) (model.User, error) {
	users, err := getUsersByBsonDocument(d)
	if err != nil {
		return model.User{}, err
	}

	if len(users) <= 0 {
		return model.User{}, errors.New("no entry found")
	}

	return users[0], nil
}

func GetUsers() ([]model.User, error) {
	return getUsersByBsonDocument(bson.D{})
}

func GetUserById(id primitive.ObjectID) (model.User, error) {
	return getUserByBsonDocument(bson.D{{"_id", id}})
}

func GetUserByKeycloakId(id string) (model.User, error) {
	return getUserByBsonDocument(bson.D{{"keycloak_id", id}})
}

func RemoveUserById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddUser(user model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := collection.InsertOne(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return GetUserById(r.InsertedID.(primitive.ObjectID))
}

func UpdateUser(user model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.D{{"_id", user.Identifier}}, user)
	if err != nil {
		return model.User{}, err
	}

	return GetUserById(user.Identifier)
}
