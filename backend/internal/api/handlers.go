package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jai-dewani/trash-taste-search/internal/db"
	"github.com/jai-dewani/trash-taste-search/internal/models"
	"github.com/jai-dewani/trash-taste-search/internal/search"
)

// Handler contains all API handlers
type Handler struct {
	db            *db.DB
	searchService *search.Service
}

// NewHandler creates a new Handler instance
func NewHandler(database *db.DB) *Handler {
	return &Handler{
		db:            database,
		searchService: search.New(database),
	}
}

// HealthCheck handles GET /api/health
func (h *Handler) HealthCheck(c *gin.Context) {
	dbStatus := "ok"
	if err := h.db.Ping(); err != nil {
		dbStatus = "error: " + err.Error()
	}

	c.JSON(http.StatusOK, models.HealthResponse{
		Status:   "ok",
		Database: dbStatus,
	})
}

// Search handles GET /api/search?q=query
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_query",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	// Get optional limit parameter (default 50)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 200 {
				limit = 200 // Cap at 200 results
			}
		}
	}

	response, err := h.searchService.Search(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "search_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetEpisodes handles GET /api/episodes
func (h *Handler) GetEpisodes(c *gin.Context) {
	episodes, err := h.db.GetAllEpisodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// Ensure we return an empty slice instead of nil
	if episodes == nil {
		episodes = []models.Episode{}
	}

	c.JSON(http.StatusOK, models.EpisodeListResponse{
		Count:    len(episodes),
		Episodes: episodes,
	})
}

// GetEpisodeByID handles GET /api/episodes/:id
func (h *Handler) GetEpisodeByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_id",
			Message: "Episode ID is required",
		})
		return
	}

	episode, err := h.db.GetEpisodeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	if episode == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Episode not found",
		})
		return
	}

	// Get segments for this episode
	segments, err := h.db.GetSegmentsByEpisodeID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// Ensure we return an empty slice instead of nil
	if segments == nil {
		segments = []models.Segment{}
	}

	c.JSON(http.StatusOK, models.EpisodeDetailResponse{
		Episode:  *episode,
		Segments: segments,
	})
}
