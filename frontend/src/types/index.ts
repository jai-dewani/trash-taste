// Episode represents a podcast episode with metadata
export interface Episode {
  id: string;
  title: string;
  description: string;
  publishedAt: string;
  channelId?: string;
  channelTitle?: string;
  thumbnailUrl: string;
}

// Segment represents a transcript segment with timestamps
export interface Segment {
  id: number;
  episodeId: string;
  startTime: number;
  endTime: number;
  text: string;
}

// SearchResult represents a search result with episode info and matching segment
export interface SearchResult {
  episode: Episode;
  segment: Segment;
  highlight?: string;
}

// SearchResponse is the API response for search queries
export interface SearchResponse {
  query: string;
  count: number;
  results: SearchResult[];
}

// Sort options for results
export type SortOption = 'relevance' | 'date_asc' | 'date_desc';
