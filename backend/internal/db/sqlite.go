package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"

	"github.com/jai-dewani/trash-taste-search/internal/models"
)

// DB wraps the SQLite database connection
type DB struct {
	conn *sql.DB
	mu   sync.RWMutex
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// GetAllEpisodes returns all episodes from the database
func (db *DB) GetAllEpisodes() ([]models.Episode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	rows, err := db.conn.Query(`
		SELECT id, title, description, published_at, channel_id, channel_title, thumbnail_url
		FROM episodes
		ORDER BY published_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query episodes: %w", err)
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var ep models.Episode
		var channelID, channelTitle sql.NullString
		err := rows.Scan(
			&ep.ID,
			&ep.Title,
			&ep.Description,
			&ep.PublishedAt,
			&channelID,
			&channelTitle,
			&ep.ThumbnailURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan episode: %w", err)
		}
		ep.ChannelID = channelID.String
		ep.ChannelTitle = channelTitle.String
		episodes = append(episodes, ep)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating episodes: %w", err)
	}

	return episodes, nil
}

// GetEpisodeByID returns a single episode by ID
func (db *DB) GetEpisodeByID(id string) (*models.Episode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var ep models.Episode
	var channelID, channelTitle sql.NullString
	err := db.conn.QueryRow(`
		SELECT id, title, description, published_at, channel_id, channel_title, thumbnail_url
		FROM episodes
		WHERE id = ?
	`, id).Scan(
		&ep.ID,
		&ep.Title,
		&ep.Description,
		&ep.PublishedAt,
		&channelID,
		&channelTitle,
		&ep.ThumbnailURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query episode: %w", err)
	}

	ep.ChannelID = channelID.String
	ep.ChannelTitle = channelTitle.String
	return &ep, nil
}

// GetSegmentsByEpisodeID returns all segments for an episode
func (db *DB) GetSegmentsByEpisodeID(episodeID string) ([]models.Segment, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	rows, err := db.conn.Query(`
		SELECT id, episode_id, start_time, end_time, text
		FROM segments
		WHERE episode_id = ?
		ORDER BY start_time ASC
	`, episodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query segments: %w", err)
	}
	defer rows.Close()

	var segments []models.Segment
	for rows.Next() {
		var seg models.Segment
		err := rows.Scan(&seg.ID, &seg.EpisodeID, &seg.StartTime, &seg.EndTime, &seg.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to scan segment: %w", err)
		}
		segments = append(segments, seg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating segments: %w", err)
	}

	return segments, nil
}

// SearchSegments performs a full-text search on transcript segments
func (db *DB) SearchSegments(query string, limit int) ([]models.SearchResult, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}

	// Use FTS5 to search for matching segments
	rows, err := db.conn.Query(`
		SELECT 
			s.id, s.episode_id, s.start_time, s.end_time, s.text,
			e.id, e.title, e.description, e.published_at, e.channel_id, e.channel_title, e.thumbnail_url,
			highlight(segments_fts, 0, '<mark>', '</mark>') as highlighted
		FROM segments_fts
		JOIN segments s ON segments_fts.rowid = s.id
		JOIN episodes e ON s.episode_id = e.id
		WHERE segments_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search segments: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var result models.SearchResult
		var channelID, channelTitle sql.NullString
		err := rows.Scan(
			&result.Segment.ID,
			&result.Segment.EpisodeID,
			&result.Segment.StartTime,
			&result.Segment.EndTime,
			&result.Segment.Text,
			&result.Episode.ID,
			&result.Episode.Title,
			&result.Episode.Description,
			&result.Episode.PublishedAt,
			&channelID,
			&channelTitle,
			&result.Episode.ThumbnailURL,
			&result.Highlight,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		result.Episode.ChannelID = channelID.String
		result.Episode.ChannelTitle = channelTitle.String
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return results, nil
}

// GetEpisodeCount returns the total number of episodes
func (db *DB) GetEpisodeCount() (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM episodes").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count episodes: %w", err)
	}
	return count, nil
}

// GetSegmentCount returns the total number of segments
func (db *DB) GetSegmentCount() (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM segments").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count segments: %w", err)
	}
	return count, nil
}
