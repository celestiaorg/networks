#!/usr/bin/env bats

@test "dry-run filters tcp_verified and sorts by score desc" {
  out="$BATS_TEST_TMPDIR/cand_out.txt"
  run bash "$BATS_TEST_DIRNAME/fetch-filter.sh" --dry-run \
    "$BATS_TEST_DIRNAME/testdata/peers.sample.json" "$out"
  [ "$status" -eq 0 ]
  run diff -u "$BATS_TEST_DIRNAME/testdata/peers.sample.golden" "$out"
  [ "$status" -eq 0 ]
}
