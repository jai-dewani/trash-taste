import sqlite3
from typing import List, Dict, Any
import sys

def search_fts5(db_path: str, search_query: str, limit: int = 10) -> List[Dict[str, Any]]:
    """
    Query an FTS5 table with the given search query.
    
    Args:
        db_path: Path to the SQLite database file
        search_query: The search query string
        limit: Maximum number of results to return (default: 10)
    
    Returns:
        List of dictionaries containing the search results
    """
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    try:
        # Query the FTS5 table with MATCH and order by relevance (bm25)
        cursor.execute("""
            SELECT *, bm25(segments_fts) as rank
            FROM segments_fts
            WHERE segments_fts MATCH ?
            ORDER BY rank
            LIMIT ?
        """, (search_query, limit))
        
        results = [dict(row) for row in cursor.fetchall()]
        return results
    
    except sqlite3.Error as e:
        print(f"Database error: {e}")
        return []
    
    finally:
        conn.close()


if __name__ == "__main__":
    # Example usage
    db_path = "data/trash_taste.db"
    if len(sys.argv) < 2:
        print("Usage: python test_database.py <search_query>")
        sys.exit(1)
    
    query = sys.argv[1]
    count = int(sys.argv[2]) if len(sys.argv) > 2 else 10
    results = search_fts5(db_path, query, count)
    
    for result in results:
        print(result)