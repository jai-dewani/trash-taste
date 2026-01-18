import sqlite3
import json
import os
from pathlib import Path

def create_database(db_path: str) -> sqlite3.Connection:
    """Create SQLite database with FTS5 table for transcript search.
    This function sets up a SQLite database optimized for searching podcast/video
    transcripts. It creates:
    1. **episodes table**: Stores metadata about each episode (title, description,
        publish date, channel info, thumbnail).
    2. **segments table**: Stores individual transcript segments with timestamps,
        allowing precise navigation to specific moments in episodes.
    3. **segments_fts (FTS5 virtual table)**: A full-text search index using trigram
        tokenization, enabling fast substring matching and fuzzy search on transcript
        text.
    4. **Sync triggers**: Automatically keeps the FTS index synchronized when
        segments are inserted, updated, or deleted.
    5. **Index on episode_id**: Speeds up joins between segments and episodes.
    Use Case:
         This is designed for a "Trash Taste" podcast search application where users
         can search through episode transcripts and find exact moments where specific
         topics are discussed. The trigram tokenizer allows partial word matching
         (e.g., searching "prog" would match "programming").
    Args:
         db_path: Path where the SQLite database file will be created or opened.
    Returns:
         sqlite3.Connection: An open connection to the configured database.
    Example:
         >>> conn = create_database("transcripts.db")
         >>> # Now ready to insert episodes and segments for searching
    """
    """Create SQLite database with FTS5 table for transcript search."""
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Create episodes table to store episode metadata
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS episodes (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT,
            published_at TEXT,
            channel_id TEXT,
            channel_title TEXT,
            thumbnail_url TEXT
        )
    ''')
    
    # Create segments table to store transcript segments
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS segments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            episode_id TEXT NOT NULL,
            start_time REAL NOT NULL,
            end_time REAL NOT NULL,
            text TEXT NOT NULL,
            FOREIGN KEY (episode_id) REFERENCES episodes(id)
        )
    ''')
    
    # Create FTS5 virtual table for full-text search on segments
    cursor.execute('''
        CREATE VIRTUAL TABLE IF NOT EXISTS segments_fts USING fts5(
            text,
            content='segments',
            content_rowid='id',
            tokenize='trigram'
        )
    ''')
    
    # Create triggers to keep FTS table in sync
    cursor.execute('''
        CREATE TRIGGER IF NOT EXISTS segments_ai AFTER INSERT ON segments BEGIN
            INSERT INTO segments_fts(rowid, text) VALUES (new.id, new.text);
        END
    ''')
    
    cursor.execute('''
        CREATE TRIGGER IF NOT EXISTS segments_ad AFTER DELETE ON segments BEGIN
            INSERT INTO segments_fts(segments_fts, rowid, text) VALUES('delete', old.id, old.text);
        END
    ''')
    
    cursor.execute('''
        CREATE TRIGGER IF NOT EXISTS segments_au AFTER UPDATE ON segments BEGIN
            INSERT INTO segments_fts(segments_fts, rowid, text) VALUES('delete', old.id, old.text);
            INSERT INTO segments_fts(rowid, text) VALUES (new.id, new.text);
        END
    ''')
    
    # Create index on episode_id for faster joins
    cursor.execute('CREATE INDEX IF NOT EXISTS idx_segments_episode_id ON segments(episode_id)')
    
    conn.commit()
    return conn

def load_videos_metadata(videos_json_path: str) -> dict:
    """Load video metadata from videos.json."""
    with open(videos_json_path, 'r', encoding='utf-8') as f:
        videos = json.load(f)
    
    # Create a mapping from video ID to metadata
    video_map = {}
    for video in videos:
        video_id = video.get('id') or video.get('video_id')
        if video_id:
            video_map[video_id] = video
    
    return video_map

def load_transcript(transcript_path: str) -> list:
    """Load transcript segments from a JSON file."""
    with open(transcript_path, 'r', encoding='utf-8') as f:
        return json.load(f)

def populate_database(conn: sqlite3.Connection, transcripts_dir: str, videos_metadata: dict):
    """Populate the database with episodes and transcript segments."""
    cursor = conn.cursor()
    
    transcripts_path = Path(transcripts_dir)
    transcript_files = list(transcripts_path.glob('*.json'))
    
    print(f"Found {len(transcript_files)} transcript files")
    
    for transcript_file in transcript_files:
        episode_id = transcript_file.stem  # filename without extension
        
        # Get episode metadata
        metadata = videos_metadata.get(episode_id, {})
        
        # Insert episode
        cursor.execute('''
            INSERT OR REPLACE INTO episodes (id, title, description, published_at, channel_id, channel_title, thumbnail_url)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        ''', (
            episode_id,
            metadata.get('title', f'Episode {episode_id}'),
            metadata.get('description', ''),
            metadata.get('published_at') or metadata.get('publishedAt', ''),
            metadata.get('channel_id') or metadata.get('channelId', ''),
            metadata.get('channel_title') or metadata.get('channelTitle', ''),
            metadata.get('thumbnail_url') or metadata.get('thumbnail', {}).get('url', '') if isinstance(metadata.get('thumbnail'), dict) else metadata.get('thumbnail', '')
        ))
        
        # Load and insert transcript segments
        try:
            segments = load_transcript(transcript_file)
            
            # Delete existing segments for this episode (in case of re-run)
            cursor.execute('DELETE FROM segments WHERE episode_id = ?', (episode_id,))
            
            for segment in segments:
                # Handle different possible transcript formats
                text = segment.get('text', '')
                start_time = segment.get('start', 0)
                duration = segment.get('duration', 0)
                end_time = start_time + duration
                
                if text.strip():  # Only insert non-empty segments
                    cursor.execute('''
                        INSERT INTO segments (episode_id, start_time, end_time, text)
                        VALUES (?, ?, ?, ?)
                    ''', (episode_id, start_time, end_time, text))
            
            print(f"Processed: {episode_id} - {len(segments)} segments")
            
        except Exception as e:
            print(f"Error processing {transcript_file}: {e}")
    
    conn.commit()
    print("Database population complete!")

def main():
    # Define paths
    base_dir = Path(__file__).parent.parent
    db_path = base_dir / 'scripts' / 'data' / 'trash_taste.db'
    transcripts_dir = base_dir / 'scripts' / 'data' / 'transcripts'
    videos_json_path = base_dir / 'scripts' / 'data' / 'videos.json'
    
    # Ensure data directory exists
    db_path.parent.mkdir(parents=True, exist_ok=True)
    
    # Remove existing database if present (for clean rebuild)
    if db_path.exists():
        os.remove(db_path)
        print(f"Removed existing database: {db_path}")
    
    # Create database
    print(f"Creating database at: {db_path}")
    conn = create_database(str(db_path))
    
    # Load video metadata
    print(f"Loading video metadata from: {videos_json_path}")
    videos_metadata = load_videos_metadata(str(videos_json_path))
    print(f"Loaded metadata for {len(videos_metadata)} videos")
    
    # Populate database
    print(f"Loading transcripts from: {transcripts_dir}")
    populate_database(conn, str(transcripts_dir), videos_metadata)
    
    # Print summary
    cursor = conn.cursor()
    cursor.execute('SELECT COUNT(*) FROM episodes')
    episode_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM segments')
    segment_count = cursor.fetchone()[0]
    
    print(f"\nDatabase Summary:")
    print(f"  Episodes: {episode_count}")
    print(f"  Segments: {segment_count}")
    
    # Test FTS search
    print("\nTesting FTS search for 'anime'...")
    cursor.execute('''
        SELECT e.title, s.text, s.start_time
        FROM segments_fts fts
        JOIN segments s ON fts.rowid = s.id
        JOIN episodes e ON s.episode_id = e.id
        WHERE segments_fts MATCH 'bread is bad'
        LIMIT 5
    ''')
    results = cursor.fetchall()
    for title, text, start_time in results:
        print(f"  [{title}] @ {start_time:.1f}s: {text[:80]}...")
    
    conn.close()
    print(f"\nDatabase saved to: {db_path}")

if __name__ == '__main__':
    main()