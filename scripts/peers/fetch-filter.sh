#!/usr/bin/env bash
# Fetch the Cumulo peer feed (or read a local fixture) and emit candidate
# "id@host:port" lines: tcp_verified peers, sorted by score descending.
#
# The feed is only a candidate source. Every line it emits is independently
# re-verified by scripts/peers/verifypeer before anything is committed.
#
# Usage:
#   fetch-filter.sh <cumulo_url> <out_file> <expected_chain_id>
#   fetch-filter.sh --dry-run <fixture_json> <out_file> <expected_chain_id>
set -euo pipefail

FILTER='.peers | sort_by(.score) | reverse | .[] | select(.tcp_verified == true) | .connection_string'

if [[ "${1:-}" == "--dry-run" ]]; then
  src="${2:?fixture path required}"
  out="${3:?out file required}"
  expected_chain_id="${4:?expected chain id required}"
  feed="$src"
else
  url="${1:?cumulo url required}"
  out="${2:?out file required}"
  expected_chain_id="${3:?expected chain id required}"
  feed="$(mktemp)"
  trap 'rm -f "$feed"' EXIT
  curl -fsSL "$url" -o "$feed"
fi

# The feed declares which chain it describes. Cumulo repointed its "testnet"
# endpoint from mocha-4 to mocha-5 without changing the URL; without this check
# that silently rewrites the wrong network's peers.txt.
actual_chain_id="$(jq -r '.chain_id // empty' "$feed")"
if [[ "$actual_chain_id" != "$expected_chain_id" ]]; then
  echo "error: feed chain_id is '${actual_chain_id:-<missing>}', expected '${expected_chain_id}'" >&2
  exit 1
fi

jq -r "$FILTER" "$feed" > "$out"

if [[ ! -s "$out" ]]; then
  echo "error: no candidates after filtering to tcp_verified peers" >&2
  exit 1
fi
