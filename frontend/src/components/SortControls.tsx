import type { SortOption } from '../types';

interface SortControlsProps {
  value: SortOption;
  onChange: (value: SortOption) => void;
  resultCount: number;
}

export function SortControls({ value, onChange, resultCount }: SortControlsProps) {
  return (
    <div className="sort-controls">
      <span className="result-count">
        {resultCount} result{resultCount !== 1 ? 's' : ''} found
      </span>
      <div className="sort-options">
        <label htmlFor="sort-select">Sort by:</label>
        <select
          id="sort-select"
          value={value}
          onChange={(e) => onChange(e.target.value as SortOption)}
          className="sort-select"
        >
          <option value="relevance">Relevance</option>
          <option value="date_desc">Newest First</option>
          <option value="date_asc">Oldest First</option>
        </select>
      </div>
    </div>
  );
}
