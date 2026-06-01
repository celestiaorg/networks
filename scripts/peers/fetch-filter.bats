#!/usr/bin/env bats

@test "dry-run filters tcp_verified and sorts by score desc" {
  run bash "$BATS_TEST_DIRNAME/fetch-filter.sh" --dry-run \
    "$BATS_TEST_DIRNAME/testdata/peers.sample.json" /tmp/cand_out.txt
  [ "$status" -eq 0 ]
  run diff -u "$BATS_TEST_DIRNAME/testdata/peers.sample.golden" /tmp/cand_out.txt
  [ "$status" -eq 0 ]
}
