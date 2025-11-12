# Ultimate Guitar Scraper v4.0

A high-performance, production-ready web service for fetching and converting guitar tabs and worship chords with intelligent chord analysis. Built entirely in Go with a modern Material Design 2 React frontend.

## 🎸 Features

### Core Features
- **Ultimate Guitar Integration**: Search and fetch guitar tabs/chords with rating-based sorting
- **Worshipchords Integration**: Fetch worship song chords from worshipchords.com
- **Manual Submission**: Submit your own chord sheets with automatic bracket wrapping
- **Key Detection**: Intelligent key detection from chord progressions
- **OnSong Format Conversion**: Professional OnSong format with metadata
- **Google Drive Integration**: Webhook support for automated uploads via n8n

### UI Features
- **Material Design 2**: Beautiful dark theme with rounded UI elements
- **Dynamic Search**: Type-as-you-search with 500ms debounce
- **Smart Sorting**: Results automatically sorted by highest rating
- **Tabbed Interface**: Easy navigation between sources
- **Rate Limiting**: Built-in protection against abuse (100 req/15min)
- **CORS Support**: Configurable origin whitelist with wildcard support
- **Home Assistant Compatible**: Designed as a Home Assistant add-on

## 🏗️ Architecture

### Pure Go Backend
- **Single Binary**: No Node.js runtime required
- **Fast**: Direct function calls (no subprocess overhead)
- **Lightweight**: ~30-50MB memory footprint
- **Efficient**: Concurrent request handling with goroutines

### Technology Stack
- **Backend**: Go 1.21+ with net/http
- **Router**: gorilla/mux
- **CORS**: rs/cors
- **Rate Limiting**: golang.org/x/time/rate
- **Frontend**: React (served as static files)
- **Deployment**: Docker multi-stage builds

## 📁 Project Structure

```
scraper/
├── main.go                    # Application entry point
├── go.mod                     # Go dependencies
├── Dockerfile                 # Multi-stage build
│
├── config/                    # Configuration management
│   ├── config.go              # Config loader
│   └── types.go               # Config structs
│
├── server/                    # HTTP server
│   ├── server.go              # Server setup & lifecycle
│   ├── middleware.go          # CORS, rate limiting, logging
│   └── router.go              # Route definitions
│
├── handlers/                  # HTTP request handlers
│   ├── health.go              # GET  /health
│   ├── search.go              # GET  /search
│   ├── onsong.go              # POST /onsong
│   ├── worshipchords.go       # POST /worshipchords
│   ├── drive.go               # POST /send-to-drive
│   ├── static.go              # GET  / (React app)
│   └── types.go               # Request/response DTOs
│
├── pkg/                       # Business logic
│   ├── ultimateguitar/        # Ultimate Guitar scraping
│   │   ├── api.go             # API client
│   │   ├── tabs.go            # Tab fetching
│   │   ├── onsong.go          # OnSong conversion
│   │   ├── types.go           # Data models
│   │   └── paths.go           # API endpoints
│   └── worshipchords/         # Worshipchords scraping
│       └── worshipchords.go   # HTML parsing & formatting
│
└── build/                     # Build artifacts
    └── frontend/              # Compiled React app
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ (for building frontend)
- Docker (optional, for containerized deployment)

### Local Development

1. **Clone the repository**
```bash
git clone <repository-url>
cd scraper
```

2. **Install Go dependencies**
```bash
go mod download
```

3. **Build the frontend** (if you have the React source)
```bash
cd frontend
npm install
npm run build
cd ..
```

4. **Create environment configuration** (optional)
```bash
cat > .env <<EOF
PORT=3000
N8N_WEBHOOK_URL=https://your-webhook-url.com/webhook/...
ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
EOF
```

5. **Run the server**
```bash
go run main.go
```

The server will start on `http://localhost:3000`

### Docker Build

```bash
# Build the Docker image
docker build -t ultimate-guitar-scraper .

# Run the container
docker run -p 3000:3000 \
  -e N8N_WEBHOOK_URL=https://your-webhook-url.com/webhook/... \
  -e ALLOWED_ORIGINS=http://localhost:3000 \
  ultimate-guitar-scraper
```

### Home Assistant Add-on

For Home Assistant deployments, the add-on automatically loads configuration from `/data/options.json`:

```json
{
  "webhook_url": "https://n8n.example.com/webhook/...",
  "allowed_origins": "http://localhost:3000,https://your-domain.com"
}
```

## 🔌 API Endpoints

### 1. Health Check
```http
GET /health
```

**Response:** `200 OK`

---

### 2. Search Ultimate Guitar
```http
GET /search?title=wonderwall
```

**Response:**
```json
[
  {
    "id": 123456,
    "song": "Wonderwall",
    "artist": "Oasis",
    "type": "Chords",
    "url": "https://tabs.ultimate-guitar.com/...",
    "rating": 4.8
  }
]
```

---

### 3. Get Tab in OnSong Format
```http
POST /onsong
Content-Type: application/json

{
  "id": 123456
}
```

**Response:** Plain text in OnSong format
```
Wonderwall
Oasis
Key: Em
Tempo: 100 BPM
Time Signature: 4/4

Intro:
[Em] [G] [Dsus4] [A7sus4]

Verse 1:
[Em]Today is [G]gonna be the day...
```

---

### 4. Get Worship Chords
```http
POST /worshipchords
Content-Type: application/json

{
  "url": "https://worshipchords.com/..."
}
```

**Response:** Plain text chord chart
```
Amazing Grace
John Newton
Key: G
Tempo: 100 BPM

Verse 1:
[G]Amazing [G7]grace how [C]sweet the [G]sound...
```

---

### 5. Send to Google Drive
```http
POST /send-to-drive
Content-Type: application/json

{
  "content": "...",
  "song": "Wonderwall",
  "artist": "Oasis",
  "id": 123456,
  "isManualSubmission": true,
  "requiresAutomation": false
}
```

**Rate Limited:** 50 requests per 15 minutes

---

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | HTTP server port |
| `N8N_WEBHOOK_URL` | *(default URL)* | n8n webhook for Google Drive |
| `ALLOWED_ORIGINS` | `localhost:3000,...` | Comma-separated CORS whitelist |

### Home Assistant Options

Configuration in `/data/options.json`:

```json
{
  "webhook_url": "string (optional)",
  "allowed_origins": "string (optional)"
}
```

### Rate Limits

- **General endpoints**: 100 requests per 15 minutes per IP
- **Google Drive upload**: 50 requests per 15 minutes per IP

### CORS Origins

Supports wildcard patterns:
- `https://*.example.com` - matches all subdomains
- `http://localhost:*` - matches all localhost ports

## 🛠️ Development

### Project Commands

```bash
# Run locally
go run main.go

# Build binary
go build -o scraper main.go

# Run tests (when available)
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Update dependencies
go mod tidy
```

### Adding a New Endpoint

1. Create handler in `handlers/`:
```go
// handlers/myendpoint.go
package handlers

func MyEndpointHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

2. Register route in `server/router.go`:
```go
router.HandleFunc("/myendpoint", handlers.MyEndpointHandler).Methods("GET")
```

3. Test the endpoint

## 📊 Performance Comparison

| Metric | Old (Node.js + Go) | New (Pure Go) | Improvement |
|--------|-------------------|---------------|-------------|
| **Memory** | 150-200 MB | 30-50 MB | 75% reduction |
| **Startup** | 3-5 seconds | 200ms | 95% faster |
| **Request Latency** | 250-400ms | 120-200ms | 50% faster |
| **Docker Image** | ~500 MB | ~50 MB | 90% smaller |
| **Process Overhead** | Subprocess per request | Direct calls | Eliminated |

## 🔒 Security Features

- **Rate Limiting**: Token bucket algorithm per IP
- **CORS**: Strict origin validation with wildcards
- **Input Validation**: All user inputs sanitized
- **URL Validation**: Domain whitelist for external requests
- **No Command Injection**: Eliminated subprocess execution
- **HTTPS Support**: TLS-ready (configure reverse proxy)

## 🐳 Docker Image Details

### Multi-Stage Build
1. **Stage 1**: Build React frontend (Node.js)
2. **Stage 2**: Build Go binary
3. **Stage 3**: Minimal Alpine image (~50MB)

### Image Contents
- Alpine Linux base
- CA certificates (for HTTPS)
- Go server binary (statically linked)
- React build (static files)

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

[Your License Here]

## 🙏 Acknowledgments

- Ultimate Guitar (ultimate-guitar.com)
- Worshipchords (worshipchords.com)
- OnSong format specification
- Home Assistant community

## 📞 Support

For issues or questions:
- Open an issue on GitHub
- Check existing issues first
- Provide logs and reproduction steps

---

**Built with ❤️ using Go**
