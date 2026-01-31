package models

// Episode represents a podcast episode with metadata
type Episode struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PublishedAt  string `json:"publishedAt"`
	ChannelID    string `json:"channelId,omitempty"`
	ChannelTitle string `json:"channelTitle,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

// Segment represents a transcript segment with timestamps
type Segment struct {
	ID        int64   `json:"id"`
	EpisodeID string  `json:"episodeId"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	Text      string  `json:"text"`
}

// SearchResult represents a search result with episode info and matching segment
type SearchResult struct {
	Episode   Episode `json:"episode"`
	Segment   Segment `json:"segment"`
	Highlight string  `json:"highlight,omitempty"`
}

// SearchResponse is the API response for search queries
type SearchResponse struct {
	Query   string         `json:"query"`
	Count   int            `json:"count"`
	Results []SearchResult `json:"results"`
}

// EpisodeListResponse is the API response for episode list
type EpisodeListResponse struct {
	Count    int       `json:"count"`
	Episodes []Episode `json:"episodes"`
}

// EpisodeDetailResponse is the API response for a single episode with segments
type EpisodeDetailResponse struct {
	Episode  Episode   `json:"episode"`
	Segments []Segment `json:"segments"`
}

// HealthResponse is the API response for health check
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// ErrorResponse is the API response for errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
