package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var dashboardCollection *mongo.Collection

func dashboardService(database *mongo.Database) {
	dashboardCollection = database.Collection("dashboard")
}

func getDashboardsByBsonDocument(d interface{}) ([]model.Dashboard, error) {
	var dashboards []model.Dashboard

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := dashboardCollection.Find(ctx, d)
	if err != nil {
		return []model.Dashboard{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var dashboard model.Dashboard
		cursor.Decode(&dashboard)

		for _, container := range dashboard.WidgetContainer {
			for _, widget := range container.Widgets {
				widget.Widget, _ = getWidgetById(widget.WidgetID)
			}
		}

		dashboards = append(dashboards, dashboard)
	}

	if err := cursor.Err(); err != nil {
		return []model.Dashboard{}, err
	}

	return dashboards, nil
}

func getDashboardByBsonDocument(d interface{}) (model.Dashboard, error) {
	dashboards, err := getDashboardsByBsonDocument(d)
	if err != nil {
		return model.Dashboard{}, err
	}

	if len(dashboards) <= 0 {
		return model.Dashboard{}, errors.New(entryNotFoundMessage)
	}

	return dashboards[0], nil
}

func GetDashboardById(id primitive.ObjectID) (model.Dashboard, error) {
	return getDashboardByBsonDocument(bson.D{{"_id", id}})
}

func GetDefaultDashboard(meetingState string) (model.Dashboard, error) {
	return getDashboardByBsonDocument(
		bson.M{
			"$and": []interface{}{
				bson.M{"default": true},
				bson.M{
					"$or": []interface{}{
						bson.M{"meeting_state": meetingState},
						bson.M{"meeting_state": bson.M{"$exists": false}},
					},
				},
			},
		})
}

func getUserDashboard(meetingState string, uuid uuid.UUID) (model.Dashboard, error) {
	return getDashboardByBsonDocument(
		bson.M{
			"$and": []interface{}{
				bson.M{"user": uuid.String()},
				bson.M{
					"$or": []interface{}{
						bson.M{"meeting_state": meetingState},
						bson.M{"meeting_state": bson.M{"$exists": false}},
					},
				},
			},
		})
}

// GetDashboardForUser returns the dashboard for the user, a boolean if it is the default dashboard and if occurred an error
func GetDashboardForUser(meetingState string, uuid uuid.UUID) (*model.Dashboard, bool, error) {
	var dashboard model.Dashboard
	var err error
	isDefault := false

	dashboard, err = getUserDashboard(meetingState, uuid)
	if err != nil {
		if err.Error() == entryNotFoundMessage {
			dashboard, err = GetDefaultDashboard(meetingState)
			isDefault = true
		}
	}

	if err != nil {
		return nil, false, err
	}
	return &dashboard, isDefault, nil
}

func RemoveUserDashboard(id primitive.ObjectID, uuid uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := dashboardCollection.DeleteOne(ctx, bson.D{{"user", uuid.String()}, {"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddUserDashboard(dashboard model.Dashboard, uuid uuid.UUID) (model.Dashboard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dashboard.User = uuid.String()

	for _, container := range dashboard.WidgetContainer {
		for _, widget := range container.Widgets {
			widget.WidgetID = widget.Widget.Identifier
		}
	}

	r, err := dashboardCollection.InsertOne(ctx, dashboard)
	if err != nil {
		return model.Dashboard{}, err
	}

	newDashboard, err1 := GetDashboardById(r.InsertedID.(primitive.ObjectID))
	if err1 != nil {
		return model.Dashboard{}, err1
	}
	return newDashboard, nil
}
