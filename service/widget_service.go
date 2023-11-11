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

var widgetCollection *mongo.Collection

func widgetService(database *mongo.Database) {
	widgetCollection = database.Collection("widget")
}

func getWidgetsByBsonDocument(d interface{}) ([]model.Widget, error) {
	var widgets []model.Widget

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := widgetCollection.Find(ctx, d)
	if err != nil {
		return []model.Widget{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var widget model.Widget
		cursor.Decode(&widget)
		widgets = append(widgets, widget)
	}

	if err := cursor.Err(); err != nil {
		return []model.Widget{}, err
	}

	return widgets, nil
}

func getWidgetById(id primitive.ObjectID) (model.Widget, error) {
	widgets, err := getWidgetsByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.Widget{}, err
	}

	if len(widgets) <= 0 {
		return model.Widget{}, errors.New(entryNotFoundMessage)
	}

	return widgets[0], nil
}

func GetWidgets() ([]model.Widget, error) {
	return getWidgetsByBsonDocument(bson.D{})
}

func RemoveWidget(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := widgetCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddWidget(widget model.Widget) (model.Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := widgetCollection.InsertOne(ctx, widget)
	if err != nil {
		return model.Widget{}, err
	}

	newWidget, err1 := getWidgetById(r.InsertedID.(primitive.ObjectID))
	if err1 != nil {
		return model.Widget{}, err1
	}
	return newWidget, nil
}
