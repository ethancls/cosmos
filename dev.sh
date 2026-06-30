#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
SERVER_PID=""
DASHBOARD_PID=""

is_listening() {
  lsof -nP -iTCP:"$1" -sTCP:LISTEN >/dev/null 2>&1
}

cleanup() {
  if [ -n "$SERVER_PID" ]; then
    kill "$SERVER_PID" 2>/dev/null || true
  fi
  if [ -n "$DASHBOARD_PID" ]; then
    kill "$DASHBOARD_PID" 2>/dev/null || true
  fi
}

echo "Starting Cosmos dev environment..."

if is_listening 8080; then
  echo "[1/2] cosmos-server already listening on http://localhost:8080"
else
  echo "[1/2] Starting cosmos-server..."
  cd "$ROOT_DIR/api"
  mkdir -p data
  go run ./management management --config ./management.json --datadir ./data --port 8080 --metrics-port 9091 --log-file console --disable-anonymous-metrics &
  SERVER_PID=$!
fi

if is_listening 3000; then
  echo "[2/2] cosmos-dashboard already listening on http://localhost:3000"
else
  echo "[2/2] Starting cosmos-dashboard..."
  cd "$ROOT_DIR/dashboard"
  npm install --silent 2>/dev/null
  npm run dev &
  DASHBOARD_PID=$!
fi

sleep 2

if [ -n "$SERVER_PID" ] && ! kill -0 "$SERVER_PID" 2>/dev/null; then
  echo "cosmos-server failed to start"
  cleanup
  exit 1
fi

if [ -n "$DASHBOARD_PID" ] && ! kill -0 "$DASHBOARD_PID" 2>/dev/null; then
  echo "cosmos-dashboard failed to start"
  cleanup
  exit 1
fi

echo ""
echo "Cosmos dev running:"
echo "  Dashboard: http://localhost:3000"
echo "  Server:    http://localhost:8080"
echo "  Dex IDP:   http://localhost:8080/oauth2"
echo ""
echo "Press Ctrl+C to stop"
trap "cleanup; exit" INT TERM

if [ -n "$SERVER_PID" ] || [ -n "$DASHBOARD_PID" ]; then
  wait
else
  tail -f /dev/null
fi
