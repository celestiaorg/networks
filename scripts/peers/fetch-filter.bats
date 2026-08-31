#!/usr/bin/env bats

setup() {
  SCRIPT="$BATS_TEST_DIRNAME/fetch-filter.sh"
  FIXTURE="$BATS_TEST_DIRNAME/testdata/peers.sample.json"
  OUT="$BATS_TEST_TMPDIR/cand_out.txt"
}

@test "dry-run filters tcp_verified and sorts by score desc" {
  run bash "$SCRIPT" --dry-run "$FIXTURE" "$OUT" celestia
  [ "$status" -eq 0 ]
  run diff -u "$BATS_TEST_DIRNAME/testdata/peers.sample.golden" "$OUT"
  [ "$status" -eq 0 ]
}

# The feed publishes which chain it is describing. Cumulo repointed its
# "testnet" endpoint from mocha-4 to mocha-5 without changing the URL, so an
# unchecked feed silently writes peers for the wrong network.
@test "rejects a feed whose chain_id is not the expected one" {
  run bash "$SCRIPT" --dry-run "$FIXTURE" "$OUT" mocha-5
  [ "$status" -ne 0 ]
  [[ "$output" == *"chain_id"* ]]
  [[ "$output" == *"celestia"* ]]
  [[ "$output" == *"mocha-5"* ]]
}

@test "writes no candidate file when the chain_id check fails" {
  run bash "$SCRIPT" --dry-run "$FIXTURE" "$OUT" mocha-5
  [ "$status" -ne 0 ]
  [ ! -s "$OUT" ]
}

@test "requires an expected chain id" {
  run bash "$SCRIPT" --dry-run "$FIXTURE" "$OUT"
  [ "$status" -ne 0 ]
}

@test "fails when the feed yields no candidates" {
  empty="$BATS_TEST_TMPDIR/empty.json"
  echo '{"chain_id":"celestia","peers":[]}' > "$empty"
  run bash "$SCRIPT" --dry-run "$empty" "$OUT" celestia
  [ "$status" -ne 0 ]
  [[ "$output" == *"no candidates"* ]]
}
