package infrastructure

import (
	"context"
	"go-multiple-query/pkg/xlogger"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongodbSetup() *mongo.Database {
	logger := xlogger.Logger

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(cfg.MongoDb.URI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Close the client if an error occurs during initialization
	defer func() {
		if err != nil {
			if cerr := client.Disconnect(context.Background()); cerr != nil {
				log.Printf("Failed to disconnect from MongoDB: %v", cerr)
			}
		}
	}()

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	logger.Info().Msg("Successfully connected to MongoDB")

	db := client.Database("vip-voucher-test")
	return db
}
