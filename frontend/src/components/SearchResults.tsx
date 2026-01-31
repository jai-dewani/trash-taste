import type { SearchResult } from '../types';
import { ResultCard } from './ResultCard';

interface SearchResultsProps {
  results: SearchResult[];
  isLoading: boolean;
  error: Error | null;
  query: string;
}

export function SearchResults({ results, isLoading, error, query }: SearchResultsProps) {
  if (!query || query.length < 2) {
    return (
      <div className="search-placeholder">
        <div className="placeholder-icon">🎙️</div>
        <h2>Search Trash Taste Podcast</h2>
        <p>Type at least 2 characters to search through podcast transcripts</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Searching transcripts...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error-icon">⚠️</div>
        <h3>Search Error</h3>
        <p>{error.message}</p>
        <p className="error-hint">Make sure the backend server is running</p>
      </div>
    );
  }

  if (results.length === 0) {
    return (
      <div className="no-results">
        <div className="no-results-icon">🔍</div>
        <h3>No results found</h3>
        <p>Try different keywords or check your spelling</p>
      </div>
    );
  }

  return (
    <div className="results-grid">
      {results.map((result, index) => (
        <ResultCard key={`${result.episode.id}-${result.segment.id}-${index}`} result={result} />
      ))}
    </div>
  );
}
