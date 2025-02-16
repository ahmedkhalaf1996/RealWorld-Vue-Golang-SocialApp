package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	mongoUri := "mongodb://localhost:27017"
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err == nil {
		fmt.Println("connection to DB..")
		DB = Client.Database("social")
	} else {
		fmt.Printf("error connect to db:%s\n", err.Error())
	}
}
