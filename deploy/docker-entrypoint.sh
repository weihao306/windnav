#!/bin/sh
set -eu

/app/windnav &
BACKEND_PID=$!

nginx -g "daemon off;" &
NGINX_PID=$!

trap 'kill "$BACKEND_PID" "$NGINX_PID" 2>/dev/null || true; wait "$BACKEND_PID" 2>/dev/null || true; wait "$NGINX_PID" 2>/dev/null || true' INT TERM EXIT

while true; do
    if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
        wait "$BACKEND_PID"
        BACKEND_EXIT=$?
        kill "$NGINX_PID" 2>/dev/null || true
        wait "$NGINX_PID" || true
        exit "$BACKEND_EXIT"
    fi

    if ! kill -0 "$NGINX_PID" 2>/dev/null; then
        wait "$NGINX_PID"
        NGINX_EXIT=$?
        kill "$BACKEND_PID" 2>/dev/null || true
        wait "$BACKEND_PID" || true
        exit "$NGINX_EXIT"
    fi

    sleep 1
done