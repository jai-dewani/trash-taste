package search

import (
	"strings"

	"github.com/jai-dewani/trash-taste-search/internal/db"
	"github.com/jai-dewani/trash-taste-search/internal/models"
)

// Service handles search operations
type Service struct {
	db *db.DB
}

// New creates a new search service
func New(database *db.DB) *Service {
	return &Service{db: database}
}

// Search performs a full-text search on transcripts
func (s *Service) Search(query string, limit int) (*models.SearchResponse, error) {
	// Sanitize and prepare query for FTS5
	sanitizedQuery := sanitizeQuery(query)

	if sanitizedQuery == "" {
		return &models.SearchResponse{
			Query:   query,
			Count:   0,
			Results: []models.SearchResult{},
		}, nil
	}

	results, err := s.db.SearchSegments(sanitizedQuery, limit)
	if err != nil {
		return nil, err
	}

	// Ensure we return an empty slice instead of nil
	if results == nil {
		results = []models.SearchResult{}
	}

	return &models.SearchResponse{
		Query:   query,
		Count:   len(results),
		Results: results,
	}, nil
}

// sanitizeQuery prepares a query string for FTS5
// It handles special characters and formats the query appropriately
func sanitizeQuery(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)

	if query == "" {
		return ""
	}

	// Replace + with space (URL encoding)
	query = strings.ReplaceAll(query, "+", " ")

	// For FTS5 trigram tokenizer, we can use the query directly
	// Multiple words will be searched as a phrase
	// Escape special FTS5 characters
	query = strings.ReplaceAll(query, "\"", "\"\"")

	// Wrap in quotes for phrase matching
	return "\"" + query + "\""
}
