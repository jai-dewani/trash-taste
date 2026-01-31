import { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { SearchBar, SearchResults, SortControls } from './components';
import { useSearch, useSortedResults } from './hooks/useSearch';
import { useDebounce } from './hooks/useDebounce';
import type { SortOption } from './types';
import './App.css';

const queryClient = new QueryClient();

function SearchApp() {
  const [searchQuery, setSearchQuery] = useState('');
  const [sortOption, setSortOption] = useState<SortOption>('relevance');
  const debouncedQuery = useDebounce(searchQuery, 300);

  const { data, isLoading, error } = useSearch(debouncedQuery);
  const sortedResults = useSortedResults(data?.results, sortOption);

  return (
    <div className="app">
      <header className="app-header">
        <h1 className="app-title">
          <span className="title-icon">🎙️</span>
          Trash Taste Search
        </h1>
        <p className="app-subtitle">
          Search through podcast transcripts to find specific moments
        </p>
      </header>

      <main className="app-main">
        <SearchBar
          value={searchQuery}
          onChange={setSearchQuery}
          isLoading={isLoading && debouncedQuery.length >= 2}
        />

        {data && data.results.length > 0 && (
          <SortControls
            value={sortOption}
            onChange={setSortOption}
            resultCount={data.count}
          />
        )}

        <SearchResults
          results={sortedResults}
          isLoading={isLoading && debouncedQuery.length >= 2}
          error={error}
          query={debouncedQuery}
        />
      </main>

      <footer className="app-footer">
        <p>
          Search through Trash Taste podcast episodes • Built with ❤️
        </p>
      </footer>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <SearchApp />
    </QueryClientProvider>
  );
}

export default App;
