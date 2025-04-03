package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"os"

	"notes-app/internal/core/service"
	"notes-app/internal/infrastructure/adapters/persistence"
	"notes-app/internal/infrastructure/adapters/web"
	"notes-app/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && os.Getenv("GO_ENV") != "production" {
		log.Println("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB client
	client, err := database.NewMongoDBClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	mongoDatabase := database.GetMongoDBDatabase(client, "")            // El nombre se lee de env vars
	mongoCollection := database.GetMongoDBCollection(mongoDatabase, "") // El nombre se lee de env vars

	noteRepo := persistence.NewMongoNoteRepository(client, mongoDatabase.Name(), mongoCollection.Name())
	noteSvc := service.NewNoteService(noteRepo)
	noteHandler := web.NewNoteHandler(noteSvc)

	router := gin.Default()

	router.POST("/notes", noteHandler.CreateNote)
	router.GET("/notes/:id", noteHandler.GetNote)
	router.PUT("/notes/:id", noteHandler.UpdateNote)
	router.DELETE("/notes/:id", noteHandler.DeleteNote)
	router.GET("/notes", noteHandler.GetNotes)     // Opcional
	router.DELETE("/notes", noteHandler.DeleteAll) // Opcional

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
