# Stage 1: Build React Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy frontend source
COPY frontend/ .

# Build the React app
RUN npm run build

# Stage 2: Build Go Binary
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (this includes all local packages)
# Version 4.1.2 - Fix fmt.Errorf usage in onsong.go
COPY . .

# Build the Go binary
# All local packages (config/, server/, handlers/, pkg/) are in the same module
# so they will be recognized without needing go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o chord-scraper .

# Stage 3: Final Image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the Go binary from builder
COPY --from=go-builder /app/chord-scraper .

# Copy the React build from frontend builder
COPY --from=frontend-builder /app/frontend/build ./build/frontend

# Expose port 3000
EXPOSE 3000

# Run the server
CMD ["./chord-scraper"]
