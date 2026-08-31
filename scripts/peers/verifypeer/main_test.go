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

// Writing before checking the gate truncates peers.txt even on a failed run.
func TestWriteVerified_BelowMinLeavesExistingFileIntact(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")
	original := "a@1:1\nb@2:2\nc@3:3\nd@4:4\ne@5:5\n"
	writeFile(t, out, original)

	err := writeVerified(out, []string{"a@1:1", "b@2:2"}, 5, 0)
	if err == nil {
		t.Fatal("writeVerified succeeded, want an error for 2 verified < min 5")
	}
	got, readErr := os.ReadFile(out)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != original {
		t.Fatalf("peers.txt was modified by a failed run:\ngot  %q\nwant %q", got, original)
	}
}

func TestWriteVerified_SortsAndDedupes(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")

	if err := writeVerified(out, []string{"c@3:3", "a@1:1", "c@3:3", "b@2:2"}, 3, 0); err != nil {
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

// MIN_PEERS alone does not stop a curated 71-peer list collapsing to 6: the
// absolute floor is met while 92% of the list silently disappears.
func TestWriteVerified_RejectsCollapseRelativeToPrevious(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")
	var previous string
	for i := 0; i < 20; i++ {
		previous += string(rune('a'+i)) + "@1:1\n"
	}
	writeFile(t, out, previous)

	// 6 verified clears min=5 but retains only 30% of a 20-peer list.
	err := writeVerified(out, []string{"a@1:1", "b@1:1", "c@1:1", "d@1:1", "e@1:1", "f@1:1"}, 5, 50)
	if err == nil {
		t.Fatal("writeVerified2 succeeded, want an error for a 70% shrink")
	}
	got, readErr := os.ReadFile(out)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != previous {
		t.Fatal("peers.txt was modified by a run that failed the retention gate")
	}
}

func TestWriteVerified_AllowsModestShrink(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")
	writeFile(t, out, "a@1:1\nb@1:1\nc@1:1\nd@1:1\n")

	// 3 of 4 retained (75%) clears a 50% floor.
	if err := writeVerified(out, []string{"a@1:1", "b@1:1", "c@1:1"}, 3, 50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteVerified_GrowthIsAlwaysAllowed(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")
	writeFile(t, out, "a@1:1\n")

	if err := writeVerified(out, []string{"a@1:1", "b@1:1", "c@1:1"}, 1, 50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteVerified_FirstRunHasNoPreviousToCompare(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "peers.txt")

	if err := writeVerified(out, []string{"a@1:1", "b@1:1"}, 2, 50); err != nil {
		t.Fatalf("unexpected error on first run: %v", err)
	}
}
