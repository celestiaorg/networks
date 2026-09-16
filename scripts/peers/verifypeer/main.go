package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

func verifiedLines(results []Result) []string {
	var out []string
	for _, r := range results {
		if r.Status == StatusVerified {
			out = append(out, r.Line)
		}
	}
	return out
}

func mismatches(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Status == StatusIDMismatch {
			out = append(out, r)
		}
	}
	return out
}

func wrongChain(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Status == StatusWrongChain {
			out = append(out, r)
		}
	}
	return out
}

// readPeerLines reads non-empty, whitespace-trimmed lines from path.
func readPeerLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, sc.Err()
}

// dedupe removes duplicates, keeping the first occurrence and input order.
func dedupe(lines []string) []string {
	seen := make(map[string]struct{}, len(lines))
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if _, ok := seen[l]; ok {
			continue
		}
		seen[l] = struct{}{}
		out = append(out, l)
	}
	return out
}

// loadCandidates builds the set of peers to verify: the fetched candidate list
// plus, when merge is set, the peers already committed to outPath.
//
// Merging matters. The upstream feed is a candidate source, not an authority on
// which peers this repo publishes; without the merge every run replaces the
// curated list with whatever the feed happened to return that morning. A
// missing outPath is the first-run case and is not an error.
func loadCandidates(inPath, outPath string, merge bool) ([]string, error) {
	candidates, err := readPeerLines(inPath)
	if err != nil {
		return nil, err
	}
	if merge {
		existing, err := readPeerLines(outPath)
		switch {
		case err == nil:
			candidates = append(candidates, existing...)
		case errors.Is(err, fs.ErrNotExist):
			fmt.Fprintf(os.Stderr, "note: %s does not exist yet, nothing to merge\n", outPath)
		default:
			return nil, err
		}
	}
	return dedupe(candidates), nil
}

// writeVerified sorts, deduplicates and writes lines to path. Sorting keeps the
// output stable regardless of dial timing, so a rerun that finds the same peers
// produces no diff.
//
// Verifying nothing at all is a malfunction rather than a judgment call for a
// reviewer, so it is refused: writing the empty result would open a PR that
// deletes the entire peer list. Every other outcome is left to the human
// reviewing the generated PR.
func writeVerified(path string, lines []string) error {
	lines = dedupe(lines)
	sort.Strings(lines)
	if len(lines) == 0 {
		return fmt.Errorf("no peers verified; leaving %s unchanged", path)
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

// verifyAll runs VerifyPeer over candidates with bounded concurrency,
// preserving input order in the returned slice.
func verifyAll(candidates []string, chainID string, concurrency int, timeout time.Duration) []Result {
	results := make([]Result, len(candidates))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, line := range candidates {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, line string) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = VerifyPeer(line, chainID, timeout)
		}(i, line)
	}
	wg.Wait()
	return results
}

func main() {
	in := flag.String("in", "", "input file of id@host:port candidate lines")
	out := flag.String("out", "", "output file for verified peer lines")
	chainID := flag.String("chain-id", "", "expected chain ID, e.g. celestia or mocha-5 (required)")
	merge := flag.Bool("merge", true, "also verify the peers already in -out, keeping the ones that still pass")
	report := flag.String("report", "", "output file for JSON verification report")
	concurrency := flag.Int("concurrency", 16, "max concurrent dials")
	timeoutSec := flag.Int("timeout", 5, "per-peer dial+handshake timeout in seconds")
	flag.Parse()

	if *in == "" || *out == "" || *chainID == "" {
		fmt.Fprintln(os.Stderr, "usage: verifypeer -in candidates.txt -out peers.txt -chain-id celestia [-merge=false] [-report report.json] [-concurrency N] [-timeout SECONDS]")
		os.Exit(2)
	}

	candidates, err := loadCandidates(*in, *out, *merge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read candidates: %v\n", err)
		os.Exit(1)
	}

	results := verifyAll(candidates, *chainID, *concurrency, time.Duration(*timeoutSec)*time.Second)
	verified := verifiedLines(results)

	if *report != "" {
		b, _ := json.MarshalIndent(results, "", "  ")
		if err := os.WriteFile(*report, b, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write report: %v\n", err)
		}
	}

	mm := mismatches(results)
	wc := wrongChain(results)
	fmt.Printf("chain=%s candidates=%d verified=%d unreachable=%d handshake_failed=%d id_mismatch=%d wrong_chain=%d invalid=%d\n",
		*chainID, len(candidates), len(dedupe(verified)),
		countStatus(results, StatusUnreachable), countStatus(results, StatusHandshakeFailed),
		len(mm), len(wc), countStatus(results, StatusInvalid))
	for _, r := range mm {
		fmt.Fprintf(os.Stderr, "ID MISMATCH: claimed=%s remote=%s line=%s\n", r.ClaimedID, r.RemoteID, r.Line)
	}
	for _, r := range wc {
		fmt.Fprintf(os.Stderr, "WRONG CHAIN: want=%s got=%s line=%s\n", *chainID, r.Network, r.Line)
	}

	if err := writeVerified(*out, verified); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
}

func countStatus(results []Result, s Status) int {
	n := 0
	for _, r := range results {
		if r.Status == s {
			n++
		}
	}
	return n
}
