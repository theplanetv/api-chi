FROM golang:1.23.7-bookworm

WORKDIR /api-chi

# Copy package files and download
COPY go.mod go.sum ./
RUN go mod download

# Install goose database migration
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY main.go ./
COPY start.sh ./
COPY wait-for-it.sh ./
RUN chmod +x ./start.sh
RUN chmod +x ./wait-for-it.sh

# Copy migrations
COPY migrations ./migrations

# Build and run
RUN go build .

ENTRYPOINT ["./start.sh"]
