#!/usr/bin/env bash
# Configure GitHub Actions secrets and variables for the staging environment.
#
# This script never prints secret values and never passes secret values as
# command-line arguments. It requires GitHub CLI authentication with permission
# to administer repository Actions secrets, variables, and environments.

set -Eeuo pipefail

ENVIRONMENT="staging"
REPO=""

usage() {
  cat <<'USAGE'
Usage:
  scripts/configure-staging-github-actions.sh [--repo OWNER/REPOSITORY]

The script creates or updates the GitHub Actions environment "staging" and
configures the values consumed by .github/workflows/deploy-staging.yml.

Required environment secrets:
  STAGING_HOST
  STAGING_SSH_USER
  STAGING_SSH_PRIVATE_KEY
  STAGING_SSH_KNOWN_HOSTS
  STAGING_GHCR_USERNAME
  STAGING_GHCR_TOKEN

Required environment variables:
  STAGING_APP_DIR
  STAGING_HEALTHCHECK_URL

Optional environment variable:
  STAGING_PUBLIC_URL
USAGE
}

fail() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

while (($# > 0)); do
  case "$1" in
    --repo)
      (($# >= 2)) || fail "--repo requires OWNER/REPOSITORY."
      REPO="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      fail "Unknown argument: $1"
      ;;
  esac
done

command -v gh >/dev/null 2>&1 || fail "GitHub CLI (gh) is not installed."

gh auth status >/dev/null 2>&1 || fail "GitHub CLI is not authenticated. Run: gh auth login"

if [[ -z "$REPO" ]]; then
  REPO="$(gh repo view --json nameWithOwner --jq '.nameWithOwner')" || fail "Could not determine the current repository. Use --repo OWNER/REPOSITORY."
fi

[[ "$REPO" =~ ^[^/]+/[^/]+$ ]] || fail "Repository must use OWNER/REPOSITORY format."

printf 'Repository: %s\n' "$REPO"
printf 'Environment: %s\n\n' "$ENVIRONMENT"

# PUT is idempotent: it creates the environment when absent and preserves it
# when already present. Protection rules should be configured separately in the
# GitHub UI or through the REST API according to the deployment policy.
printf 'Ensuring GitHub environment exists...\n'
gh api --method PUT "repos/$REPO/environments/$ENVIRONMENT" >/dev/null

set_secret_from_prompt() {
  local name="$1"
  local prompt="$2"
  local value
  local tmp

  read -r -s -p "$prompt: " value
  printf '\n'
  [[ -n "$value" ]] || fail "$name cannot be empty."

  tmp="$(mktemp)"
  chmod 600 "$tmp"
  trap 'rm -f "$tmp"' RETURN
  printf '%s' "$value" >"$tmp"
  gh secret set "$name" --repo "$REPO" --env "$ENVIRONMENT" <"$tmp" >/dev/null
  rm -f "$tmp"
  trap - RETURN
  printf 'Configured secret: %s\n' "$name"
}

set_secret_from_file() {
  local name="$1"
  local prompt="$2"
  local path

  read -r -p "$prompt: " path
  [[ -n "$path" ]] || fail "$name file path cannot be empty."
  [[ -f "$path" ]] || fail "File not found for $name: $path"
  [[ -s "$path" ]] || fail "File is empty for $name: $path"

  gh secret set "$name" --repo "$REPO" --env "$ENVIRONMENT" <"$path" >/dev/null
  printf 'Configured secret: %s\n' "$name"
}

set_variable_from_prompt() {
  local name="$1"
  local prompt="$2"
  local value

  read -r -p "$prompt: " value
  [[ -n "$value" ]] || fail "$name cannot be empty."
  gh variable set "$name" --repo "$REPO" --env "$ENVIRONMENT" --body "$value" >/dev/null
  printf 'Configured variable: %s\n' "$name"
}

printf '%s\n' 'Enter staging connection settings. Secret values remain hidden.'
read -r -p 'STAGING_HOST: ' staging_host
[[ -n "$staging_host" ]] || fail 'STAGING_HOST cannot be empty.'
gh secret set STAGING_HOST --repo "$REPO" --env "$ENVIRONMENT" --body "$staging_host" >/dev/null
printf '%s\n' 'Configured secret: STAGING_HOST'

read -r -p 'STAGING_SSH_USER: ' staging_ssh_user
[[ -n "$staging_ssh_user" ]] || fail 'STAGING_SSH_USER cannot be empty.'
gh secret set STAGING_SSH_USER --repo "$REPO" --env "$ENVIRONMENT" --body "$staging_ssh_user" >/dev/null
printf '%s\n' 'Configured secret: STAGING_SSH_USER'

set_secret_from_file STAGING_SSH_PRIVATE_KEY 'Path to STAGING_SSH_PRIVATE_KEY file'
set_secret_from_file STAGING_SSH_KNOWN_HOSTS 'Path to STAGING_SSH_KNOWN_HOSTS file'
set_secret_from_prompt STAGING_GHCR_USERNAME 'STAGING_GHCR_USERNAME'
set_secret_from_prompt STAGING_GHCR_TOKEN 'STAGING_GHCR_TOKEN'

printf '%s\n' 'Enter non-secret staging variables.'
set_variable_from_prompt STAGING_APP_DIR 'STAGING_APP_DIR'
set_variable_from_prompt STAGING_HEALTHCHECK_URL 'STAGING_HEALTHCHECK_URL'

read -r -p 'STAGING_PUBLIC_URL (optional, press Enter to skip): ' staging_public_url
if [[ -n "$staging_public_url" ]]; then
  gh variable set STAGING_PUBLIC_URL --repo "$REPO" --env "$ENVIRONMENT" --body "$staging_public_url" >/dev/null
  printf '%s\n' 'Configured variable: STAGING_PUBLIC_URL'
else
  printf '%s\n' 'Skipped optional variable: STAGING_PUBLIC_URL'
fi

printf '\nConfiguration completed for %s/%s.\n' "$REPO" "$ENVIRONMENT"
printf '%s\n' 'Secret values were not displayed. Re-run the staging workflow after verifying the host and health endpoint.'
