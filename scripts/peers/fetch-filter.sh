#!/usr/bin/env bash
# Fetch the Cumulo peer feed (or read a local fixture) and emit candidate
# "id@host:port" lines: tcp_verified peers, sorted by score descending.
#
# Usage:
#   fetch-filter.sh <cumulo_url> <out_file>
#   fetch-filter.sh --dry-run <fixture_json> <out_file>
set -euo pipefail

FILTER='.peers | sort_by(.score) | reverse | .[] | select(.tcp_verified == true) | .connection_string'

if [[ "${1:-}" == "--dry-run" ]]; then
  fixture="${2:?fixture path required}"
  out="${3:?out file required}"
  jq -r "$FILTER" "$fixture" > "$out"
  exit 0
fi

url="${1:?cumulo url required}"
out="${2:?out file required}"
curl -fsSL "$url" | jq -r "$FILTER" > "$out"
