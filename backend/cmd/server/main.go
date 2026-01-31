// Entry point for the server application

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jai-dewani/trash-taste-search/internal/api"
	"github.com/jai-dewani/trash-taste-search/internal/db"
)

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "../data/trash_taste.db"
	}

	// Get server port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	log.Printf("Connecting to database: %s", dbPath)
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Get episode and segment counts for logging
	episodeCount, _ := database.GetEpisodeCount()
	segmentCount, _ := database.GetSegmentCount()
	log.Printf("Database loaded: %d episodes, %d segments", episodeCount, segmentCount)

	// Setup router
	router := api.SetupRouter(database)

	// Setup CORS
	corsHandler := api.NewCORSHandler()

	// Create server with CORS middleware
	handler := corsHandler.Handler(router)

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Printf("API endpoints:")
	log.Printf("  GET /api/health           - Health check")
	log.Printf("  GET /api/search?q=<query> - Search transcripts")
	log.Printf("  GET /api/episodes         - List all episodes")
	log.Printf("  GET /api/episodes/:id     - Get episode details")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
