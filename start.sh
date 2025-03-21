#!/bin/sh
set -e

# Run migration until success
./wait-for-it.sh ${POSTGRES_PUBLIC_HOST}:${POSTGRES_PORT} -t 30

# Run migration
echo "Running migration"
goose up

# Start server
echo "Starting server..."
exec ./api-chi
