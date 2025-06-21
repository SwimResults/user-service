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

var reportCollection *mongo.Collection

func reportService(database *mongo.Database) {
	reportCollection = database.Collection("report")
}

func getReportsByBsonDocument(d interface{}) ([]model.UserReport, error) {
	var reports []model.UserReport

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := reportCollection.Find(ctx, d)
	if err != nil {
		return []model.UserReport{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var report model.UserReport
		cursor.Decode(&report)
		reports = append(reports, report)
	}

	if err := cursor.Err(); err != nil {
		return []model.UserReport{}, err
	}

	return reports, nil
}

func getReportById(id primitive.ObjectID) (model.UserReport, error) {
	reports, err := getReportsByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.UserReport{}, err
	}

	if len(reports) <= 0 {
		return model.UserReport{}, errors.New(entryNotFoundMessage)
	}

	return reports[0], nil
}

func GetReports() ([]model.UserReport, error) {
	return getReportsByBsonDocument(bson.D{})
}

func GetReportsByMeeting(meeting string) ([]model.UserReport, error) {
	return getReportsByBsonDocument(bson.D{{"meeting", meeting}})
}

func RemoveReport(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := reportCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddReport(report model.UserReport) (model.UserReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	report.AddedAt = time.Now()

	r, err := reportCollection.InsertOne(ctx, report)
	if err != nil {
		return model.UserReport{}, err
	}

	newReport, err1 := getReportById(r.InsertedID.(primitive.ObjectID))
	if err1 != nil {
		return model.UserReport{}, err1
	}
	return newReport, nil
}

func UpdateReport(report model.UserReport) (model.UserReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := reportCollection.ReplaceOne(ctx, bson.D{{"_id", report.Identifier}}, report)
	if err != nil {
		return model.UserReport{}, err
	}

	newReport, err1 := getReportById(report.Identifier)
	if err1 != nil {
		return model.UserReport{}, err1
	}
	return newReport, nil
}

func ToggleReportAcknowledged(id primitive.ObjectID) (model.UserReport, error) {
	report, err := getReportById(id)
	if err != nil {
		return model.UserReport{}, err
	}

	report.Acknowledged = !report.Acknowledged

	return UpdateReport(report)
}

func ToggleReportComplete(id primitive.ObjectID) (model.UserReport, error) {
	report, err := getReportById(id)
	if err != nil {
		return model.UserReport{}, err
	}

	report.Completed = !report.Completed

	return UpdateReport(report)
}
