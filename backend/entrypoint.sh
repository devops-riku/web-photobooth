#!/bin/sh
set -e

ENV_PATH=".env"
TEMPLATE_PATH=".env.template"

if [ ! -f "$ENV_PATH" ]; then
  if [ -f "$TEMPLATE_PATH" ]; then
    echo "[entrypoint] .env not found, creating from .env.template"
    cp "$TEMPLATE_PATH" "$ENV_PATH"
  else
    echo "[entrypoint] No .env.template found, using system env vars"
  fi
else
  echo "[entrypoint] .env already exists"
fi

exec "$@"
