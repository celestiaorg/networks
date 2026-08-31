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

// writeVerified sorts, deduplicates and writes lines to path, but only once two
// gates pass. Both are checked before the write, so a failed run leaves the
// existing peers.txt untouched rather than truncating it.
//
//   - min: an absolute floor on the number of verified peers.
//   - minRetainPct: the new list must be at least this percentage of the list
//     already at path. The absolute floor alone does not catch a curated
//     71-peer list collapsing to 6 because upstream had a bad morning; growth
//     and the first run (no file yet) are always allowed. 0 disables the gate.
func writeVerified(path string, lines []string, min, minRetainPct int) error {
	lines = dedupe(lines)
	sort.Strings(lines) // deterministic output independent of dial timing
	if len(lines) < min {
		return fmt.Errorf("only %d verified peers, need >= %d; leaving %s unchanged", len(lines), min, path)
	}
	if minRetainPct > 0 {
		previous, err := readPeerLines(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if n := len(dedupe(previous)); n > 0 && len(lines)*100 < n*minRetainPct {
			return fmt.Errorf("verified %d peers but %s currently has %d (%d%%, below the %d%% floor); "+
				"leaving it unchanged \u2014 rerun or lower -min-retain-pct if this shrink is intended",
				len(lines), path, n, len(lines)*100/n, minRetainPct)
		}
	}
	content := []byte{}
	if len(lines) > 0 {
		content = []byte(strings.Join(lines, "\n") + "\n")
	}
	return os.WriteFile(path, content, 0o644)
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
	min := flag.Int("min", 5, "minimum verified peers required (else exit 1)")
	concurrency := flag.Int("concurrency", 16, "max concurrent dials")
	minRetainPct := flag.Int("min-retain-pct", 50, "fail if the new list is below this percentage of the existing -out list (0 disables)")
	timeoutSec := flag.Int("timeout", 5, "per-peer dial+handshake timeout in seconds")
	flag.Parse()

	if *in == "" || *out == "" || *chainID == "" {
		fmt.Fprintln(os.Stderr, "usage: verifypeer -in candidates.txt -out peers.txt -chain-id celestia [-merge=false] [-report report.json] [-min N] [-min-retain-pct N] [-concurrency N] [-timeout SECONDS]")
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

	if err := writeVerified(*out, verified, *min, *minRetainPct); err != nil {
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
