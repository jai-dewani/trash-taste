# 🎙️ Trash Taste Search

A full-stack application that allows you to search through **Trash Taste podcast** episode transcripts and jump directly to specific moments. Find that one clip where they talked about *that thing* — instantly.

![Tech Stack](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-07405E?style=for-the-badge&logo=sqlite&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-646CFF?style=for-the-badge&logo=vite&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)

## ✨ Features

- 🔍 **Full-Text Search** — Search through 220+ episode transcripts with blazing-fast FTS5 trigram search
- ⏱️ **Timestamp Links** — Click any result to jump directly to that moment on YouTube
- 📅 **Sort Options** — Sort results by relevance or date (newest/oldest first)
- 🎨 **Modern UI** — Clean, responsive interface built with React and Tailwind CSS
- ⚡ **Fast & Efficient** — Go backend with SQLite for sub-millisecond search queries
- 🔄 **Real-time Search** — Debounced search with React Query for smooth UX

## 🖼️ How It Works

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│   React Frontend │────▶│    Go Backend    │────▶│  SQLite + FTS5   │
│   (Vite + TS)    │◀────│   (Gin Router)   │◀────│  (Trigram Index) │
└──────────────────┘     └──────────────────┘     └──────────────────┘
         │                                                  │
         │                                                  │
         ▼                                                  ▼
   User searches for                               220+ episodes with
   "Connor's hot take"                             timestamped segments
```

1. **User types a search query** in the frontend
2. **Frontend debounces** the input and sends a request to the API
3. **Backend performs FTS5 search** using trigram tokenization for substring matching
4. **Results include episode metadata + matching transcript segments** with timestamps
5. **User clicks a result** → Opens YouTube at the exact timestamp

## 🛠️ Tech Stack

### Backend
- **[Go](https://golang.org/)** — Fast, compiled language for the API server
- **[Gin](https://gin-gonic.com/)** — High-performance HTTP router
- **[SQLite](https://sqlite.org/)** with **FTS5** — Full-text search with trigram tokenization
- **[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)** — Pure Go SQLite driver (no CGO)

### Frontend
- **[React 19](https://react.dev/)** — UI library
- **[TypeScript](https://www.typescriptlang.org/)** — Type-safe JavaScript
- **[Vite](https://vitejs.dev/)** (Rolldown) — Next-gen frontend tooling
- **[TanStack Query](https://tanstack.com/query)** — Server state management
- **[Tailwind CSS](https://tailwindcss.com/)** — Utility-first CSS framework
- **[Axios](https://axios-http.com/)** — HTTP client

### Scripts (Data Pipeline)
- **Python 3** — For fetching and processing data
- **YouTube Data API v3** — Fetch video metadata
- **youtube-transcript-api** — Download auto-generated transcripts

## 📁 Project Structure

```
trash-taste-search/
├── backend/                    # Go API server
│   ├── cmd/server/main.go      # Entry point
│   └── internal/
│       ├── api/                # HTTP handlers & routing
│       ├── db/                 # SQLite database operations
│       ├── models/             # Data structures
│       └── search/             # Search service logic
│
├── frontend/                   # React application
│   ├── src/
│   │   ├── api/                # API client
│   │   ├── components/         # React components
│   │   ├── hooks/              # Custom hooks (useSearch, useDebounce)
│   │   └── types/              # TypeScript interfaces
│   └── package.json
│
├── data/                       # Data storage
│   ├── videos.json             # Episode metadata (16k+ lines)
│   ├── transcripts/            # 220+ transcript JSON files
│   └── trash_taste.db          # SQLite database (generated)
│
└── scripts/                    # Python data pipeline
    ├── fetch_videos.py         # Fetch video metadata from YouTube API
    ├── download_transcript.py  # Download transcripts
    └── generate_database.py    # Build SQLite database with FTS5
```

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+**
- **Node.js 20+** (with npm/pnpm)
- **Python 3.10+** (for data scripts)
- **YouTube Data API v3 key** (for fetching new data)

### 1. Clone the Repository

```bash
git clone https://github.com/jai-dewani/trash-taste-search.git
cd trash-taste-search
```

### 2. Set Up the Database

If the database doesn't exist, generate it from transcripts:

```bash
cd scripts
pip install -r requirements.txt
python generate_database.py
```

This creates `data/trash_taste.db` with:
- **episodes table** — Video metadata
- **segments table** — Timestamped transcript segments
- **segments_fts** — FTS5 virtual table with trigram tokenization

### 3. Start the Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

The server starts on `http://localhost:8080` with these endpoints:

| Endpoint | Description |
|----------|-------------|
| `GET /api/health` | Health check |
| `GET /api/search?q=<query>&limit=50` | Search transcripts |
| `GET /api/episodes` | List all episodes |
| `GET /api/episodes/:id` | Get episode details with segments |

### 4. Start the Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend starts on `http://localhost:5173`

## 📊 Data Pipeline

To fetch fresh data from YouTube:

### 1. Set Up Environment

Create a `.env` file in the `scripts/` directory:

```env
YOUTUBE_DATA_API_V3=your_api_key_here
```

### 2. Fetch Video Metadata

```bash
python fetch_videos.py
```

This fetches all videos from the Trash Taste channel and saves them to `data/videos.json`.

### 3. Download Transcripts

```bash
python download_transcript.py
```

Downloads auto-generated transcripts for each video to `data/transcripts/`.

### 4. Generate Database

```bash
python generate_database.py
```

Creates/updates the SQLite database with FTS5 full-text search.

## 🔍 Search Features

### Trigram Tokenization

The search uses SQLite FTS5 with **trigram tokenization**, which means:

- ✅ Substring matching: Searching "prog" finds "programming"
- ✅ Case-insensitive search
- ✅ Phrase matching with quoted strings
- ✅ Fast prefix/suffix matching

### Search Examples

| Query | Finds |
|-------|-------|
| `connor cycling` | Moments about Connor's cycling journey |
| `Joey anime` | Joey's anime takes |
| `Garnt Thailand` | When Garnt talks about Thailand |
| `hot take` | All the hot takes |

## 🔧 Configuration

### Backend Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_PATH` | `../data/trash_taste.db` | Path to SQLite database |
| `PORT` | `8080` | Server port |

### Frontend Environment Variables

Configure in `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080
```

## 📈 API Response Examples

### Search Response

```json
{
  "query": "connor cycling",
  "count": 42,
  "results": [
    {
      "episode": {
        "id": "abc123",
        "title": "We Became Professional Athletes | Trash Taste #200",
        "thumbnailUrl": "https://i.ytimg.com/vi/abc123/hqdefault.jpg",
        "publishedAt": "2024-01-15T20:00:00Z"
      },
      "segment": {
        "id": 12345,
        "episodeId": "abc123",
        "startTime": 1234.5,
        "endTime": 1240.2,
        "text": "So Connor, tell us about your cycling journey..."
      },
      "highlight": "So <mark>Connor</mark>, tell us about your <mark>cycling</mark> journey..."
    }
  ]
}
```

## 🤝 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is for educational purposes. Trash Taste is a podcast by Joey (The Anime Man), Connor (CDawgVA), and Garnt (Gigguk).

## 🙏 Acknowledgments

- **[Trash Taste](https://www.youtube.com/@TrashTaste)** — The best podcast on the internet
- **Joey, Connor, and Garnt** — For hundreds of hours of amazing content
- **Mudan** — For the legendary edits

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/jai-dewani">Jai Dewani</a>
</p>

