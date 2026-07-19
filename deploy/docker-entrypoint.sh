#!/bin/sh
set -eu

log() {
  printf '%s %s\n' "[entrypoint]" "$*"
}

require_file() {
  if [ ! -f "$1" ]; then
    log "missing file: $1"
    exit 1
  fi
}

require_exec() {
  if [ ! -x "$1" ]; then
    log "not executable: $1"
    exit 1
  fi
}

start_backend() {
  log "starting backend"
  /app/windnav &
  BACKEND_PID=$!
}

start_nginx() {
  log "starting nginx"
  nginx -g "daemon off;" &
  NGINX_PID=$!
}

stop_all() {
  if [ -n "${BACKEND_PID:-}" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
    kill "$BACKEND_PID" 2>/dev/null || true
  fi
  if [ -n "${NGINX_PID:-}" ] && kill -0 "$NGINX_PID" 2>/dev/null; then
    kill "$NGINX_PID" 2>/dev/null || true
  fi
}

wait_pid() {
  pid=$1
  if wait "$pid"; then
    return 0
  fi
  return $?
}

BACKEND_PID=""
NGINX_PID=""

require_exec /app/windnav
require_file /etc/nginx/conf.d/default.conf
require_file /usr/share/nginx/html/index.html
require_file /usr/share/nginx/html/assets
mkdir -p /app/data

log "self-check passed"
start_backend
start_nginx
trap 'stop_all' INT TERM EXIT

while :; do
  if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
    wait_pid "$BACKEND_PID"
    code=$?
    log "backend exited: $code"
    stop_all
    exit "$code"
  fi

  if ! kill -0 "$NGINX_PID" 2>/dev/null; then
    wait_pid "$NGINX_PID"
    code=$?
    log "nginx exited: $code"
    stop_all
    exit "$code"
  fi

  sleep 1
done
