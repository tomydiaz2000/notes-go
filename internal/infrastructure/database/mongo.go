package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBClient crea y devuelve un nuevo cliente de MongoDB.
func NewMongoDBClient(ctx context.Context) (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Reemplaza con tu URI de MongoDB
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		if closeErr := client.Disconnect(ctx); closeErr != nil {
			log.Printf("Error disconnecting during ping failure: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

// GetMongoDBDatabase obtiene la instancia de la base de datos de MongoDB.
func GetMongoDBDatabase(client *mongo.Client, dbName string) *mongo.Database {
	if dbName == "" {
		dbName = os.Getenv("MONGO_DB_NAME")
		if dbName == "" {
			dbName = "mi_app_notas" // Valor por defecto
		}
	}
	return client.Database(dbName)
}

// GetMongoDBCollection obtiene la instancia de la colección de MongoDB.
func GetMongoDBCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	if collectionName == "" {
		collectionName = os.Getenv("MONGO_COLLECTION_NAME")
		if collectionName == "" {
			collectionName = "notes" // Valor por defecto
		}
	}
	return db.Collection(collectionName)
}
