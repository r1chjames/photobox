#!/bin/sh
set -e

# Start the Go API in the background
/app/api &

# Start nginx in the foreground (keeps container alive)
exec nginx -g "daemon off;"
