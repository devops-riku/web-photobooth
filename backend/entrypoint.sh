#!/bin/sh
set -e

ENV_FILE=".env"
TEMPLATE_FILE=".env.template"

echo "[entrypoint] Starting backend bootstrap"

# -------------------------------------------------
# 1. Create .env from template if missing
# -------------------------------------------------
if [ ! -f "$ENV_FILE" ]; then
  if [ -f "$TEMPLATE_FILE" ]; then
    echo "[entrypoint] .env not found, copying from .env.template"
    cp "$TEMPLATE_FILE" "$ENV_FILE"
  else
    echo "[entrypoint] ERROR: .env.template not found"
    exit 1
  fi
else
  echo "[entrypoint] .env already exists"
fi

# -------------------------------------------------
# 2. Generate DATABASE_URL if empty (SAFE)
# -------------------------------------------------
if ! grep -q "^DATABASE_URL=.*[^[:space:]]" "$ENV_FILE"; then
  echo "[entrypoint] DATABASE_URL empty, generating"

  DB_USER="${DB_USER:-photouser}"
  DB_PASSWORD="${DB_PASSWORD:-photopass}"
  DB_NAME="${DB_NAME:-photobooth}"

  DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}?sslmode=disable"

  sed -i "/^DATABASE_URL=/d" "$ENV_FILE"
  echo "DATABASE_URL=${DATABASE_URL}" >> "$ENV_FILE"
fi

# -------------------------------------------------
# 3. Generate JWT_SECRET if empty (SAFE)
# -------------------------------------------------
JWT_VALUE=$(grep "^JWT_SECRET=" "$ENV_FILE" | cut -d= -f2-)

if [ -z "$JWT_VALUE" ]; then
  echo "[entrypoint] JWT_SECRET empty, generating"

  JWT_SECRET=$(openssl rand -hex 32)

  sed -i "/^JWT_SECRET=/d" "$ENV_FILE"
  echo "JWT_SECRET=${JWT_SECRET}" >> "$ENV_FILE"
fi

# -------------------------------------------------
# 4. Validate REQUIRED external secrets (DO Spaces)
# -------------------------------------------------
REQUIRED_VARS="
DO_SPACES_KEY
DO_SPACES_SECRET
DO_SPACES_ENDPOINT
DO_SPACES_REGION
DO_SPACES_BUCKET
"

for VAR in $REQUIRED_VARS; do
  VALUE=$(grep "^$VAR=" "$ENV_FILE" | cut -d= -f2-)
  if [ -z "$VALUE" ]; then
    echo "[entrypoint] ERROR: Required variable $VAR is not set"
    echo "[entrypoint] Please set it in .env"
    exit 1
  fi
done

echo "[entrypoint] Environment validated successfully"

# -------------------------------------------------
# 5. Start application
# -------------------------------------------------
exec "$@"
