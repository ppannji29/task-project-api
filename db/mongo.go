package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	Client     *mongo.Client
	UserCol    *mongo.Collection
	OtpCol     *mongo.Collection
	ProfileCol *mongo.Collection
}

func InitMongo() (*MongoCollections, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	if mongoURI == "" || dbName == "" {
		log.Fatal("Missing MONGO_URI or MONGO_DB_NAME in environment")
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Ping to confirm connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	log.Println("✅ Connected to MongoDB")

	db := client.Database(dbName)

	return &MongoCollections{
		Client:     client,
		UserCol:    db.Collection("User"),
		OtpCol:     db.Collection("Otp"),
		ProfileCol: db.Collection("Profile"),
	}, nil
}
