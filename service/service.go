package service

import (
	"context"
	"github.com/swimresults/user-service/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var mongoClient *mongo.Client
var keycloak *model.Tkc

func Init(c *mongo.Client, kc *model.Tkc) {
	database := c.Database(os.Getenv("SR_USER_MONGO_DATABASE"))
	mongoClient = c
	keycloak = kc

	userService(database)
	widgetService(database)
	dashboardService(database)
}

func PingDatabase() bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	err := mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	return true
}
