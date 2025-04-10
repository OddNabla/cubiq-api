package setup

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func Init() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	MongoClient, err = mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic("Error connecting to MongoDB")
	}
	err = MongoClient.Ping(context.Background(), nil)
	if err != nil {
		panic("Error pinging MongoDB")
	}
	// defer MongoClient.Disconnect(context.Background())
	log.Println("Connected to MongoDB")
	MongoDatabase = MongoClient.Database(os.Getenv("MONGODB_DATABASE"))

}
