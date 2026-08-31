package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestVerifiedLines(t *testing.T) {
	results := []Result{
		{Line: "a@1:1", ClaimedID: "a", Status: StatusVerified},
		{Line: "b@2:2", ClaimedID: "b", Status: StatusUnreachable},
		{Line: "c@3:3", ClaimedID: "c", Status: StatusIDMismatch},
		{Line: "d@4:4", ClaimedID: "d", Status: StatusVerified},
		{Line: "e@5:5", ClaimedID: "e", Status: StatusWrongChain},
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

func TestReadPeerLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cand.txt")
	writeFile(t, path, "  a@1.2.3.4:26656  \n\n\nb@5.6.7.8:26656\n   \n")
	got, err := readPeerLines(path)
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

func TestDedupe(t *testing.T) {
	got := dedupe([]string{"b", "a", "b", "c", "a"})
	want := []string{"b", "a", "c"} // first occurrence wins, order preserved
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

// The feed is a candidate source, not the whole truth: peers already committed
// to peers.txt must be re-verified and kept if they still pass, otherwise every
// run silently deletes the curated list.
func TestLoadCandidates_MergesExistingOutput(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "candidates.txt")
	out := filepath.Join(dir, "peers.txt")
	writeFile(t, in, "a@1:1\nb@2:2\n")
	writeFile(t, out, "b@2:2\nc@3:3\n")

	got, err := loadCandidates(in, out, true)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a@1:1", "b@2:2", "c@3:3"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestLoadCandidates_MergeDisabled(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "candidates.txt")
	out := filepath.Join(dir, "peers.txt")
	writeFile(t, in, "a@1:1\n")
	writeFile(t, out, "c@3:3\n")

	got, err := loadCandidates(in, out, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "a@1:1" {
		t.Fatalf("got %v, want [a@1:1]", got)
	}
}

func TestLoadCandidates_MissingOutputIsNotAnError(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "candidates.txt")
	writeFile(t, in, "a@1:1\n")

	got, err := loadCandidates(in, filepath.Join(dir, "does-not-exist.txt"), true)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "a@1:1" {
		t.Fatalf("got %v, want [a@1:1]", got)
	}
}

func TestWriteVerified_SortsAndDedupes(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")

	if err := writeVerified(out, []string{"c@3:3", "a@1:1", "c@3:3", "b@2:2"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	want := "a@1:1\nb@2:2\nc@3:3\n"
	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestVerifyAll(t *testing.T) {
	if got := verifyAll(nil, "celestia", 4, time.Second); len(got) != 0 {
		t.Fatalf("empty input: got %d results, want 0", len(got))
	}
	// Invalid lines resolve to StatusInvalid without any network I/O.
	in := []string{"not-a-peer-0", "not-a-peer-1", "not-a-peer-2"}
	got := verifyAll(in, "celestia", 2, time.Second)
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
