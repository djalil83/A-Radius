#!/usr/bin/env bash
set -Eeuo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
JWT_TOKEN="${JWT_TOKEN:-}"
PROFILE_ID="${PROFILE_ID:-}"
VERSION="${VERSION:-1}"

if [[ -z "$JWT_TOKEN" ]]; then
  echo "JWT_TOKEN wajib diisi dengan JWT yang diterbitkan identity provider." >&2
  exit 1
fi

json_file() { mktemp; }
request() {
  local method="$1" path="$2" body="${3:-}" expected="$4" out
  out="$(json_file)"
  if [[ -n "$body" ]]; then
    status=$(curl --silent --show-error --request "$method" "$BASE_URL$path" \
      --header "Authorization: Bearer $JWT_TOKEN" \
      --header 'Content-Type: application/json' \
      --data "$body" --output "$out" --write-out '%{http_code}')
  else
    status=$(curl --silent --show-error --request "$method" "$BASE_URL$path" \
      --header "Authorization: Bearer $JWT_TOKEN" \
      --output "$out" --write-out '%{http_code}')
  fi
  echo "$method $path -> $status" >&2
  cat "$out" >&2
  echo >&2
  [[ "$status" == "$expected" ]] || { echo "expected HTTP $expected, got $status" >&2; exit 1; }
  cat "$out"
  rm -f "$out"
}

# Missing bearer token must be rejected.
status=$(curl --silent --output /dev/null --write-out '%{http_code}' "$BASE_URL/api/v1/subscription-profiles")
echo "GET without bearer -> $status"
[[ "$status" == "401" ]] || { echo "expected HTTP 401" >&2; exit 1; }

create_body='{
  "name":"Smoke Plan 20M",
  "service_type":"FTTH",
  "category":"HOME",
  "media":"FIBER",
  "color":"#2563EB",
  "description":"Temporary API smoke-test profile",
  "mikrotik_group":"default",
  "radius_group":"ftth",
  "rate_limit":"20M/20M",
  "upload_bps":20000000,
  "download_bps":20000000,
  "shared_users":1,
  "monthly_price":150000,
  "active_days":30,
  "commission_amount":0,
  "commission_type":"FIXED",
  "billing_cycle":"MONTHLY",
  "auto_isolate":true,
  "billing_note":"Delete or archive after smoke test"
}'

created=$(request POST /api/v1/subscription-profiles "$create_body" 201)
PROFILE_ID="$(jq -r '.id' <<<"$created")"
VERSION="$(jq -r '.version' <<<"$created")"
[[ -n "$PROFILE_ID" && "$PROFILE_ID" != "null" ]] || { echo "create response has no id" >&2; exit 1; }

request GET "/api/v1/subscription-profiles/$PROFILE_ID" '' 200 >/dev/null
request GET "/api/v1/subscription-profiles/$PROFILE_ID/revisions" '' 200 >/dev/null

update_body=$(jq -n --argjson version "$VERSION" '{version:$version, name:"Smoke Plan 20M Updated", service_type:"FTTH", color:"#1D4ED8", monthly_price:160000, active_days:30, commission_amount:0, commission_type:"FIXED", billing_cycle:"MONTHLY", auto_isolate:true}')
updated=$(request PATCH "/api/v1/subscription-profiles/$PROFILE_ID" "$update_body" 200)
NEXT_VERSION="$(jq -r '.version' <<<"$updated")"
[[ "$NEXT_VERSION" -gt "$VERSION" ]] || { echo "version did not increment" >&2; exit 1; }

# Reusing the old version must produce optimistic-lock conflict.
stale_body=$(jq -n --argjson version "$VERSION" '{version:$version, name:"Stale Update", service_type:"FTTH", color:"#1D4ED8", monthly_price:160000, active_days:30, commission_amount:0, commission_type:"FIXED", billing_cycle:"MONTHLY", auto_isolate:true}')
request PATCH "/api/v1/subscription-profiles/$PROFILE_ID" "$stale_body" 409 >/dev/null
request DELETE "/api/v1/subscription-profiles/$PROFILE_ID?version=$NEXT_VERSION" '' 204 >/dev/null

echo "profile API smoke test passed for $PROFILE_ID"
