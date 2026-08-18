#!/usr/bin/env bash
set -Eeuo pipefail

: "${PGHOST:?PGHOST is required}"
: "${PGPORT:=5432}"
: "${PGUSER:?PGUSER is required}"
: "${PGDATABASE:?PGDATABASE is required}"
: "${BACKUP_DIR:=/backups}"
: "${RETENTION_DAYS:=30}"
: "${BACKUP_PREFIX:=a-radius}"

umask 077
mkdir -p "$BACKUP_DIR"
lock_dir="$BACKUP_DIR/.backup.lock"
if ! mkdir "$lock_dir" 2>/dev/null; then
  echo "another backup is already running" >&2
  exit 2
fi
release_lock() { rmdir "$lock_dir" 2>/dev/null || true; }
trap release_lock EXIT

export PGPASSWORD="${PGPASSWORD:-}"
timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
base="$BACKUP_DIR/${BACKUP_PREFIX}_${PGDATABASE}_${timestamp}"
tmp_dump="${base}.dump.partial"
dump="${base}.dump"
sha="${dump}.sha256"
meta="${dump}.metadata"

cleanup() { rm -f "$tmp_dump"; }
trap cleanup EXIT

pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE"
pg_dump --format=custom --compress=6 --no-owner --no-privileges \
  --file="$tmp_dump" -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" "$PGDATABASE"
pg_restore --list "$tmp_dump" >/dev/null
mv -f "$tmp_dump" "$dump"
sha256sum "$dump" > "$sha"
{
  printf 'created_at_utc=%s\n' "$timestamp"
  printf 'database=%s\n' "$PGDATABASE"
  printf 'server=%s:%s\n' "$PGHOST" "$PGPORT"
  printf 'format=custom\n'
  printf 'sha256=%s\n' "$(cut -d' ' -f1 "$sha")"
} > "$meta"

find "$BACKUP_DIR" -maxdepth 1 -type f -name "${BACKUP_PREFIX}_${PGDATABASE}_*.dump" -mtime "+$RETENTION_DAYS" -print -delete
find "$BACKUP_DIR" -maxdepth 1 -type f \( -name "${BACKUP_PREFIX}_${PGDATABASE}_*.dump.sha256" -o -name "${BACKUP_PREFIX}_${PGDATABASE}_*.dump.metadata" \) -mtime "+$RETENTION_DAYS" -print -delete

printf 'backup_created=%s\n' "$dump"
