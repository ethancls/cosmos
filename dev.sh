#!/bin/bash
set -e
echo "Starting Cosmos dev environment..."
echo "[1/2] Building and starting cosmos-server..."
cd cosmos-server
go build -o cosmos-server . 2>/dev/null || echo "  (using existing binary)"
mkdir -p data
./cosmos-server management --config ./management.json --datadir ./data --port 8080 --log-file console --disable-anonymous-metrics &
SERVER_PID=$!
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
