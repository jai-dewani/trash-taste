import { useQuery } from '@tanstack/react-query';
import { searchApi } from '../api/searchApi';
import type { SearchResponse, SearchResult, SortOption } from '../types';
import { useMemo } from 'react';

export function useSearch(query: string, enabled: boolean = true) {
  return useQuery<SearchResponse>({
    queryKey: ['search', query],
    queryFn: () => searchApi.search(query),
    enabled: enabled && query.length >= 2,
    staleTime: 1000 * 60 * 5, // Cache for 5 minutes
    refetchOnWindowFocus: false,
  });
}

export function useSortedResults(
  results: SearchResult[] | undefined,
  sortOption: SortOption
): SearchResult[] {
  return useMemo(() => {
    if (!results || results.length === 0) return [];

    const sortedResults = [...results];

    switch (sortOption) {
      case 'date_asc':
        return sortedResults.sort(
          (a, b) =>
            new Date(a.episode.publishedAt).getTime() -
            new Date(b.episode.publishedAt).getTime()
        );
      case 'date_desc':
        return sortedResults.sort(
          (a, b) =>
            new Date(b.episode.publishedAt).getTime() -
            new Date(a.episode.publishedAt).getTime()
        );
      case 'relevance':
      default:
        // Keep original order (relevance from API)
        return sortedResults;
    }
  }, [results, sortOption]);
}
