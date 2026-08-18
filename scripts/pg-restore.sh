#!/usr/bin/env bash
set -Eeuo pipefail

: "${PGHOST:?PGHOST is required}"
: "${PGPORT:=5432}"
: "${PGUSER:?PGUSER is required}"
: "${PGDATABASE:?PGDATABASE is required}"
: "${BACKUP_FILE:?BACKUP_FILE is required}"
: "${CONFIRM_RESTORE:?Set CONFIRM_RESTORE=YES to continue}"

if [[ "$CONFIRM_RESTORE" != "YES" ]]; then
  echo 'refusing restore: set CONFIRM_RESTORE=YES explicitly' >&2
  exit 2
fi
if [[ ! -f "$BACKUP_FILE" ]]; then
  echo "backup file not found: $BACKUP_FILE" >&2
  exit 1
fi
if [[ ! -f "${BACKUP_FILE}.sha256" ]]; then
  echo "checksum file not found: ${BACKUP_FILE}.sha256" >&2
  exit 1
fi

umask 077
export PGPASSWORD="${PGPASSWORD:-}"
expected="$(awk '{print $1}' "${BACKUP_FILE}.sha256")"
actual="$(sha256sum "$BACKUP_FILE" | awk '{print $1}')"
if [[ -z "$expected" || "$expected" != "$actual" ]]; then
  echo "checksum mismatch for $BACKUP_FILE" >&2
  exit 1
fi

pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE"
pg_restore --list "$BACKUP_FILE" >/dev/null

if [[ "${ALLOW_DROP:-NO}" == "YES" ]]; then
  pg_restore --clean --if-exists --exit-on-error --no-owner --no-privileges \
    -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" "$BACKUP_FILE"
else
  pg_restore --exit-on-error --no-owner --no-privileges \
    -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" "$BACKUP_FILE"
fi

printf 'restore_completed=%s\n' "$BACKUP_FILE"
