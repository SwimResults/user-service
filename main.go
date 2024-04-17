package main

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/swimresults/user-service/controller"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var mongoClient *mongo.Client

func main() {
	ctx := connectDB()

	kc := model.Tkc{
		Realm:  "",
		Token:  nil,
		Client: nil,
	}
	connectKeycloak(&kc)

	service.Init(mongoClient, &kc)
	controller.Run()

	if err := mongoClient.Disconnect(ctx); err != nil {
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
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_USER_MONGO_HOST") + ":" + os.Getenv("SR_USER_MONGO_PORT") + "' as '" + os.Getenv("SR_USER_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to connect to mongo database"))
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_USER_MONGO_HOST") + ":" + os.Getenv("SR_USER_MONGO_PORT") + "' as '" + os.Getenv("SR_USER_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to reach mongo database"))
	}

	return ctx
}

func connectKeycloak(kc *model.Tkc) {

	url, p := os.LookupEnv("SR_USER_KEYCLOAK_URL")
	if !p {
		fmt.Println("no keycloak url given! Please set SR_USER_KEYCLOAK_URL.")
		return
	}

	user, p := os.LookupEnv("SR_USER_KEYCLOAK_USER")
	if !p {
		fmt.Println("no keycloak user given! Please set SR_USER_KEYCLOAK_USER.")
		return
	}

	password, p := os.LookupEnv("SR_USER_KEYCLOAK_PASSWORD")
	if !p {
		fmt.Println("no keycloak password given! Please set SR_USER_KEYCLOAK_PASSWORD.")
		return
	}

	realm, p := os.LookupEnv("SR_USER_KEYCLOAK_REALM")
	if !p {
		fmt.Println("no keycloak realm given! Please set SR_USER_KEYCLOAK_REALM.")
		return
	}

	kc.Realm = realm

	client := gocloak.NewClient(url)
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, user, password, "master")
	if err != nil {
		fmt.Println("failed when trying to connect to '" + url + "' as '" + user + "'")
		fmt.Println(fmt.Errorf("unable to reach keycloak"))
	}

	kc.Token = token
	kc.Client = client
}
