#!/bin/bash
set -e
echo "Starting Cosmos dev environment..."
echo "[1/2] Starting cosmos-server..."
cd api
mkdir -p data
go run ./management/ --datadir ./data --port 8080 --log-file console --disable-anonymous-metrics &
SERVER_PID=$!
echo "[2/2] Starting cosmos-dashboard..."
cd ../dashboard
npm install --silent 2>/dev/null
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
