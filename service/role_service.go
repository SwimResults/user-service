package service

import "go.mongodb.org/mongo-driver/mongo"

var roleCollection *mongo.Collection

func roleService(database *mongo.Database) {
	roleCollection = database.Collection("role")
}
