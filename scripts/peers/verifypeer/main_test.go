package main

import "testing"

func TestVerifiedLines(t *testing.T) {
	results := []Result{
		{Line: "a@1:1", ClaimedID: "a", Status: StatusVerified},
		{Line: "b@2:2", ClaimedID: "b", Status: StatusUnreachable},
		{Line: "c@3:3", ClaimedID: "c", Status: StatusIDMismatch},
		{Line: "d@4:4", ClaimedID: "d", Status: StatusVerified},
	}
	got := verifiedLines(results)
	want := []string{"a@1:1", "d@4:4"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestMismatches(t *testing.T) {
	results := []Result{
		{Line: "c@3:3", Status: StatusIDMismatch},
		{Line: "a@1:1", Status: StatusVerified},
	}
	if n := len(mismatches(results)); n != 1 {
		t.Fatalf("mismatches = %d, want 1", n)
	}
}
