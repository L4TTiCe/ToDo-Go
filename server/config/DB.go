package config

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Client

var ToDoItemsCollection *mongo.Collection

func constructURI() string {
	uri := os.Getenv("MONGODB_URI")

	MONGODB_USER := os.Getenv("MONGODB_USER")
	MONGODB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	MONGODB_PREFIX := os.Getenv("MONGODB_PREFIX")
	MONGODB_HOST := os.Getenv("MONGODB_HOST")
	MONGODB_PORT := os.Getenv("MONGODB_DOCKER_PORT")

	if uri == "" {
		if (MONGODB_USER != "") && (MONGODB_PASSWORD != "") && (MONGODB_PREFIX != "") && (MONGODB_HOST != "") {
			log.Println("Using environment variables for DB connection")
			uri = MONGODB_PREFIX + "://" + MONGODB_USER + ":" + MONGODB_PASSWORD + "@" + MONGODB_HOST
			if MONGODB_PORT != "" {
				uri += ":" + MONGODB_PORT
			}
		} else {
			log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
	} else {
		log.Println("Using 'MONGODB_URI' environmental variable for DB connection")
	}
	return uri
}

func ConnectMongoDB() {
	uri := constructURI()

	log.Println("Connecting to DB")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Successfully connected and pinged DB")

	DB = client
	linkCollections()
}

func CloseClientDB() {
	log.Println("Disconnecting from DB")
	if err := DB.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func linkCollections() {
	log.Println("Linking Collections...")
	ToDoItemsCollection = DB.Database("test").Collection("ToDoItems")
}
