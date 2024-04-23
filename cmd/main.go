package main

import (
	"context"
	"log"

	"github.com/elusiv0/medods_test/internal/app"
	"github.com/elusiv0/medods_test/internal/di"
	userModel "github.com/elusiv0/medods_test/internal/repo/user/model"
	mongoClient "github.com/elusiv0/medods_test/pkg/mongo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error with extract env varialbes")
	}
	ctn, err := di.InitContainer()
	if err != nil {
		log.Fatal("error with init app deps")
	}
	user1 := userModel.User{
		UUID: "d4d46a09-dc0c-4d66-8840-7424ce91db72",
		Name: "user1",
	}
	user2 := userModel.User{
		UUID: "09fd5cdf-cf73-46a2-bea5-7db7e82797f6",
		Name: "user2",
	}
	ctn.Get("mongo").(*mongoClient.MongoClient).MongoDatabase.Collection("users").InsertOne(context.Background(), user1)
	ctn.Get("mongo").(*mongoClient.MongoClient).MongoDatabase.Collection("users").InsertOne(context.Background(), user2)

	app := ctn.Get("app").(*app.App)
	if err := app.Run(); err != nil {
		log.Fatal("error with starting application")
	}
}
