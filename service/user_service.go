package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/swimresults/service-core/misc"
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var collection *mongo.Collection

const entryNotFoundMessage = "no entry found"

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
		return model.User{}, errors.New(entryNotFoundMessage)
	}

	return users[0], nil
}

func GetUsers() ([]model.User, error) {
	return getUsersByBsonDocument(bson.D{})
}

func GetUserById(id primitive.ObjectID) (model.User, error) {
	return getUserByBsonDocument(bson.D{{"_id", id}})
}

func GetUsersByIsFollowerOrMe(athleteIds []primitive.ObjectID) ([]model.User, error) {
	return getUsersByBsonDocument(bson.M{
		"$or": []interface{}{
			bson.M{"own_athlete_id": bson.M{"$in": athleteIds}},
			bson.M{"following.athlete_id": bson.M{"$in": athleteIds}},
		},
	})
}

func GetUsersByIsFollower(athleteIds []primitive.ObjectID) ([]model.User, error) {
	return getUsersByBsonDocument(bson.M{
		"$or": []interface{}{
			bson.M{"following.athlete_id": bson.M{"$in": athleteIds}},
		},
	})
}

func GetUsersByIsMe(athleteIds []primitive.ObjectID) ([]model.User, error) {
	return getUsersByBsonDocument(bson.M{
		"$or": []interface{}{
			bson.M{"own_athlete_id": bson.M{"$in": athleteIds}},
		},
	})
}

// GetUserByKeycloakId gets a user by keycloak id, creates new one if not existing so far
func GetUserByKeycloakId(id uuid.UUID) (model.User, error) {
	user, err := getUserByBsonDocument(bson.D{{"keycloak_id", id.String()}})
	if err != nil {
		if err.Error() == entryNotFoundMessage {

			user.KeycloakId = id.String()

			user, err = AddUser(user)
			if err != nil {
				return model.User{}, err
			}

		} else {
			return model.User{}, err
		}
	}
	return user, nil
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

func ModifyFollowForUser(id uuid.UUID, athleteId primitive.ObjectID, follow bool) (model.User, error) {
	user, err := GetUserByKeycloakId(id)
	if err != nil {
		return model.User{}, err
	}

	for i, following := range user.Following {
		if following.AthleteId == athleteId {
			if follow {
				return user, nil
			} else {
				user.Following = append(user.Following[:i], user.Following[i+1:]...)
			}
		}
	}

	if follow {
		user.Following = append(user.Following, model.FollowingAthlete{
			AthleteId: athleteId,
			AddedAt:   time.Now(),
		})
	}

	return UpdateUser(user)
}

func ModifyMe(id uuid.UUID, athleteId primitive.ObjectID, set bool) (model.User, error) {
	user, err := GetUserByKeycloakId(id)
	if err != nil {
		return model.User{}, err
	}

	if set {
		user.OwnAthleteId = &athleteId
	} else {
		user.OwnAthleteId = &primitive.ObjectID{}
	}

	return UpdateUser(user)
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

func ModifyUserLanguage(id uuid.UUID, language string) (model.User, error) {
	user, err := GetUserByKeycloakId(id)
	if err != nil {
		return model.User{}, err
	}

	user.Settings.Language = language

	return UpdateUser(user)
}

func ModifyUserTheme(id uuid.UUID, theme string) (model.User, error) {
	user, err := GetUserByKeycloakId(id)
	if err != nil {
		return model.User{}, err
	}

	user.Settings.Theme = theme

	return UpdateUser(user)
}

func ModifyUserMeetings(id uuid.UUID, meeting string, subscribe bool) (model.User, error) {
	user, err := GetUserByKeycloakId(id)
	if err != nil {
		return model.User{}, err
	}

	if subscribe {
		user.Meetings = misc.AppendWithoutDuplicates(user.Meetings, meeting)
	} else {
		for i, m := range user.Meetings {
			if m == meeting {
				user.Meetings = append(user.Meetings[:i], user.Meetings[i+1:]...)
				break
			}
		}
	}

	return UpdateUser(user)
}
