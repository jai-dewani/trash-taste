import type { SearchResult } from '../types';

interface ResultCardProps {
  result: SearchResult;
}

function formatTime(seconds: number): string {
  const hrs = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);

  if (hrs > 0) {
    return `${hrs}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

function getYouTubeUrl(videoId: string, startTime: number): string {
  return `https://www.youtube.com/watch?v=${videoId}&t=${Math.floor(startTime)}s`;
}

export function ResultCard({ result }: ResultCardProps) {
  const { episode, segment, highlight } = result;
  const youtubeUrl = getYouTubeUrl(episode.id, segment.startTime);

  return (
    <a
      href={youtubeUrl}
      target="_blank"
      rel="noopener noreferrer"
      className="result-card"
    >
      <div className="result-thumbnail">
        <img
          src={episode.thumbnailUrl}
          alt={episode.title}
          loading="lazy"
        />
        <span className="timestamp-badge">
          {formatTime(segment.startTime)}
        </span>
      </div>
      <div className="result-content">
        <h3 className="result-title">{episode.title}</h3>
        <p className="result-date">{formatDate(episode.publishedAt)}</p>
        <p className="result-text">
          {highlight ? (
            <span dangerouslySetInnerHTML={{ __html: highlight }} />
          ) : (
            segment.text
          )}
        </p>
        <div className="result-meta">
          <span className="time-range">
            {formatTime(segment.startTime)} - {formatTime(segment.endTime)}
          </span>
        </div>
      </div>
    </a>
  );
}
