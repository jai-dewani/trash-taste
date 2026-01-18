package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"github.com/mattn/go-sqlite3"
)

func getDBPath(t *testing.T) string {
	dbPath := filepath.Join("data", "trash_taste.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skipf("Database not found at %s. Run generate_database.py first.", dbPath)
	}
	return dbPath
}

func TestDatabaseExists(t *testing.T) {
	dbPath := getDBPath(t)
	_, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("Database file does not exist: %v", err)
	}
}

func TestDatabaseConnection(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestEpisodesTableExists(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM episodes").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query episodes table: %v", err)
	}

	if count == 0 {
		t.Error("Episodes table is empty")
	}
	t.Logf("Found %d episodes", count)
}

func TestSegmentsTableExists(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM segments").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query segments table: %v", err)
	}

	if count == 0 {
		t.Error("Segments table is empty")
	}
	t.Logf("Found %d segments", count)
}

func TestEpisodesTableSchema(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("PRAGMA table_info(episodes)")
	if err != nil {
		t.Fatalf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id":            false,
		"title":         false,
		"description":   false,
		"published_at":  false,
		"channel_id":    false,
		"channel_title": false,
		"thumbnail_url": false,
	}

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	for col, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column %s not found in episodes table", col)
		}
	}
}

func TestSegmentsTableSchema(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("PRAGMA table_info(segments)")
	if err != nil {
		t.Fatalf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id":         false,
		"episode_id": false,
		"start_time": false,
		"end_time":   false,
		"text":       false,
	}

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	for col, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column %s not found in segments table", col)
		}
	}
}

func TestFTSTableExists(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='segments_fts'").Scan(&name)
	if err != nil {
		t.Fatalf("FTS table segments_fts does not exist: %v", err)
	}
}

func TestFTSSearch(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	query := `
		SELECT e.title, s.text, s.start_time
		FROM segments_fts fts
		JOIN segments s ON fts.rowid = s.id
		JOIN episodes e ON s.episode_id = e.id
		WHERE segments_fts MATCH 'anime'
		LIMIT 5
	`

	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("FTS search failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var title, text string
		var startTime float64
		err := rows.Scan(&title, &text, &startTime)
		if err != nil {
			t.Fatalf("Failed to scan result: %v", err)
		}
		count++
		t.Logf("Found: [%s] @ %.1fs: %s", title, startTime, truncate(text, 50))
	}

	if count == 0 {
		t.Log("No results found for 'anime' search (this may be expected depending on content)")
	}
}

func TestSegmentsHaveValidEpisodeReferences(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var orphanCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM segments s
		LEFT JOIN episodes e ON s.episode_id = e.id
		WHERE e.id IS NULL
	`).Scan(&orphanCount)
	if err != nil {
		t.Fatalf("Failed to check orphan segments: %v", err)
	}

	if orphanCount > 0 {
		t.Errorf("Found %d segments without valid episode references", orphanCount)
	}
}

func TestSegmentTimestampsAreValid(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var invalidCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM segments
		WHERE start_time < 0 OR end_time < start_time
	`).Scan(&invalidCount)
	if err != nil {
		t.Fatalf("Failed to check invalid timestamps: %v", err)
	}

	if invalidCount > 0 {
		t.Errorf("Found %d segments with invalid timestamps", invalidCount)
	}
}

func TestIndexExists(t *testing.T) {
	dbPath := getDBPath(t)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_segments_episode_id'").Scan(&name)
	if err != nil {
		t.Fatalf("Index idx_segments_episode_id does not exist: %v", err)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}