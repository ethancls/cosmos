#!/bin/bash
# Cosmos dev — hot reload for both frontend and backend
# Requires: go, node, npm

set -e

echo "Starting Cosmos dev environment..."

# Start backend (requires rebuild on Go changes)
echo "[1/2] Starting cosmos-server..."
cd cosmos-server
go build -o cosmos-server ./management/ 2>/dev/null || true
mkdir -p data
./cosmos-server management \
  --config ./management.json \
  --datadir ./data \
  --port 8080 \
  --log-file console \
  --disable-anonymous-metrics &
SERVER_PID=$!

# Start frontend (Next.js hot reload)
echo "[2/2] Starting cosmos-dashboard..."
cd ../cosmos-dashboard
npm run dev &
DASHBOARD_PID=$!

echo ""
echo "Cosmos dev running:"
echo "  Dashboard: http://localhost:3000"
echo "  Server:    http://localhost:8080"
echo "  Dex IDP:   http://localhost:8080/oauth2"
echo ""
echo "Press Ctrl+C to stop"

trap "kill $SERVER_PID $DASHBOARD_PID 2>/dev/null; exit" INT TERM
wait
