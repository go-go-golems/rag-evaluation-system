#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/02-docker-compose.mysql.yml"
DUMP_PATH="${1:-/home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2}"
DB_NAME="${TTC_MYSQL_DATABASE:-ttc}"
ROOT_PASSWORD="${TTC_MYSQL_ROOT_PASSWORD:-somewordpress}"

if [[ ! -f "$DUMP_PATH" ]]; then
  echo "dump not found: $DUMP_PATH" >&2
  exit 1
fi

cd "$SCRIPT_DIR"
echo "Starting isolated MySQL container..." >&2
docker compose -f "$COMPOSE_FILE" up -d mysql >/dev/null

echo "Waiting for MySQL readiness..." >&2
for i in $(seq 1 90); do
  if docker compose -f "$COMPOSE_FILE" exec -T mysql mysqladmin ping -uroot -p"$ROOT_PASSWORD" --silent >/dev/null 2>&1; then
    break
  fi
  if [[ "$i" == "90" ]]; then
    echo "MySQL did not become ready" >&2
    docker compose -f "$COMPOSE_FILE" logs --tail=80 mysql >&2 || true
    exit 1
  fi
  sleep 2
done

echo "Resetting database $DB_NAME..." >&2
docker compose -f "$COMPOSE_FILE" exec -T mysql mysql -uroot -p"$ROOT_PASSWORD" <<SQL
DROP DATABASE IF EXISTS \`$DB_NAME\`;
CREATE DATABASE \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SQL

echo "Importing $(du -h "$DUMP_PATH" | awk '{print $1}') compressed dump. This can take several minutes..." >&2
bzcat "$DUMP_PATH" | docker compose -f "$COMPOSE_FILE" exec -T mysql mysql \
  -uroot -p"$ROOT_PASSWORD" \
  --default-character-set=utf8mb4 \
  --max_allowed_packet=512M \
  "$DB_NAME"

echo "Import complete. Summary counts:" >&2
docker compose -f "$COMPOSE_FILE" exec -T mysql mysql -uroot -p"$ROOT_PASSWORD" -N -B "$DB_NAME" <<'SQL'
SELECT post_type, post_status, COUNT(*) FROM wp_posts GROUP BY post_type, post_status ORDER BY post_type, post_status;
SQL
