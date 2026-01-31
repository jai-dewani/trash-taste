package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"

	"github.com/jai-dewani/trash-taste-search/internal/db"
)

// SetupRouter creates and configures the Gin router with all API routes
func SetupRouter(database *db.DB) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add recovery middleware to recover from panics
	router.Use(gin.Recovery())

	// Add logging middleware
	router.Use(gin.Logger())

	// Create handler with database
	handler := NewHandler(database)

	// API routes group
	api := router.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", handler.HealthCheck)

		// Search endpoint
		api.GET("/search", handler.Search)

		// Episodes endpoints
		api.GET("/episodes", handler.GetEpisodes)
		api.GET("/episodes/:id", handler.GetEpisodeByID)
	}

	return router
}

// NewCORSHandler creates a CORS handler with appropriate settings
func NewCORSHandler() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	})
}
