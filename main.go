package main

import (
	"context"
	"fmt"
	"github.com/swimresults/user-service/apns"
	"github.com/swimresults/user-service/controller"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var client *mongo.Client

func main() {
	ctx := connectDB()
	apns.Init()
	service.Init(client)
	controller.Run()

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func connectDB() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	var uri = "mongodb://"
	if os.Getenv("SR_USER_MONGO_USERNAME") != "" {
		uri += os.Getenv("SR_USER_MONGO_USERNAME") + ":" + os.Getenv("SR_USER_MONGO_PASSWORD") + "@"
	}
	uri += os.Getenv("SR_USER_MONGO_HOST") + ":" + os.Getenv("SR_USER_MONGO_PORT")
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_USER_MONGO_HOST") + ":" + os.Getenv("SR_USER_MONGO_PORT") + "' as '" + os.Getenv("SR_USER_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to connect to mongo database"))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_USER_MONGO_HOST") + ":" + os.Getenv("SR_USER_MONGO_PORT") + "' as '" + os.Getenv("SR_USER_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to reach mongo database"))
	}

	return ctx
}
