package database

import (
	"context"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	// Load the env file
	err := godotenv.Load(".env")
	if(err != nil){
		log.Fatal("Error loading .env file")
	}

	mongoUrl := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if(err != nil){
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if(err != nil){
		log.Fatal(err)
	}

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(Client *mongo.Client, CollectionName string) *mongo.Collection{
	collection :=Client.Database("cluster0").Collection(CollectionName)
	return collection
}