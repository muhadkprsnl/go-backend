// package config

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // ConnectMongoDB establishes a connection to MongoDB Atlas using MONGODB_URI from environment.
// func ConnectMongoDB() (*mongo.Client, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Get URI from environment variable
// 	mongoURI := os.Getenv("MONGODB_URI")
// 	if mongoURI == "" {
// 		log.Fatal("❌ MONGODB_URI environment variable not set")
// 	}

// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ping to verify connection
// 	if err = client.Ping(ctx, nil); err != nil {
// 		return nil, err
// 	}

//		log.Println("✅ Connected to MongoDB Atlas")
//		return client, nil
//	}

package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() (*mongo.Client, error) {
	// Load environment from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found or failed to load")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGODB_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to MongoDB")
	return client, nil
}
