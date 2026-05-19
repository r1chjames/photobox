#!/bin/sh
set -e

# Start API with a restart watchdog
while true; do
    /app/api
    echo "[watchdog] API exited, restarting in 2s..." >&2
    sleep 2
done &
API_WATCHDOG_PID=$!

# Cleanup handler — kills the watchdog loop (and any running API) on container stop
trap "kill $API_WATCHDOG_PID 2>/dev/null; kill %1 2>/dev/null" EXIT

# Wait for API to be ready before starting nginx (prevents 502 flood)
for i in $(seq 1 30); do
    if wget -qO- http://localhost:8080/api/health 2>/dev/null; then
        echo "[entrypoint] API is ready" >&2
        break
    fi
    echo "[entrypoint] waiting for API... ($i/30)" >&2
    sleep 1
done

# Start nginx in the foreground (keeps container alive)
exec nginx -g "daemon off;"
