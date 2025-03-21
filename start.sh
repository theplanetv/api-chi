#!/bin/sh
set -e

# Run migration until success
./wait-for-it.sh ${POSTGRES_PUBLIC_HOST}:${POSTGRES_PORT} -t 30

# Run migration
echo "Running migration"
goose up

# Run test
go test -v cmd/services/database.go cmd/services/database_test.go
go test -v cmd/services/auth.go cmd/services/auth_test.go
go test -v cmd/services/database.go cmd/services/blogtag.go cmd/services/blogtag_test.go
go test -v cmd/services/database.go cmd/services/blogtag.go cmd/services/blogpost.go cmd/services/blogpost_test.go
go test -v cmd/routes/auth.go cmd/routes/auth_test.go
go test -v cmd/routes/auth.go cmd/routes/blogtag.go cmd/routes/blogtag_test.go
go test -v cmd/routes/auth.go cmd/routes/blogpost.go cmd/routes/blogpost_test.go

# Start server
echo "Starting server..."
exec ./api-chi
