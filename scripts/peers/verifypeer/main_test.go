package main

import (
	"os"
	"testing"
	"time"
)

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
	got := mismatches(results)
	if len(got) != 1 {
		t.Fatalf("mismatches = %d, want 1", len(got))
	}
	if got[0].Line != "c@3:3" || got[0].Status != StatusIDMismatch {
		t.Fatalf("got %+v, want the id_mismatch entry", got[0])
	}
}

func TestReadCandidates(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/cand.txt"
	content := "  a@1.2.3.4:26656  \n\n\nb@5.6.7.8:26656\n   \n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := readCandidates(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a@1.2.3.4:26656", "b@5.6.7.8:26656"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestVerifyAll(t *testing.T) {
	if got := verifyAll(nil, 4, time.Second); len(got) != 0 {
		t.Fatalf("empty input: got %d results, want 0", len(got))
	}
	// Invalid lines resolve to StatusInvalid without any network I/O.
	in := []string{"not-a-peer-0", "not-a-peer-1", "not-a-peer-2"}
	got := verifyAll(in, 2, time.Second)
	if len(got) != len(in) {
		t.Fatalf("got %d results, want %d", len(got), len(in))
	}
	for i := range in {
		if got[i].Line != in[i] {
			t.Fatalf("order not preserved at %d: got %q want %q", i, got[i].Line, in[i])
		}
		if got[i].Status != StatusInvalid {
			t.Fatalf("line %q: status %q, want invalid", in[i], got[i].Status)
		}
	}
}
