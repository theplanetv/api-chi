# Builder stage
FROM golang:1.23.7-alpine AS builder

WORKDIR /api-chi

# Install bash and dependencies
RUN apk add --no-cache bash

# Copy package files and download
COPY go.mod go.sum ./
RUN go mod download

# Install goose database migration (only needed in builder)
RUN go install -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY main.go ./

# Build the application
RUN go build -o api-chi .

# Copy scripts (needed in both stages)
COPY start.sh wait-for-it.sh ./
RUN chmod +x ./start.sh ./wait-for-it.sh

# Runner stage
FROM alpine:latest

WORKDIR /api-chi

# Install bash and any other runtime dependencies
RUN apk add --no-cache bash

# Copy necessary files from builder
COPY --from=builder /api-chi/api-chi .
COPY --from=builder /api-chi/start.sh .
COPY --from=builder /api-chi/wait-for-it.sh .

# Copy migrations (if needed at runtime)
COPY migrations ./migrations

# Optional: copy goose binary if you need it at runtime
COPY --from=builder /go/bin/goose /usr/local/bin/goose

ENTRYPOINT ["./start.sh"]
