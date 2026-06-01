package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
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

func readCandidates(path string) ([]string, error) {
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

// verifyAll runs VerifyPeer over candidates with bounded concurrency,
// preserving input order in the returned slice.
func verifyAll(candidates []string, concurrency int, timeout time.Duration) []Result {
	results := make([]Result, len(candidates))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, line := range candidates {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, line string) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = VerifyPeer(line, timeout)
		}(i, line)
	}
	wg.Wait()
	return results
}

func main() {
	in := flag.String("in", "", "input file of id@host:port candidate lines")
	out := flag.String("out", "", "output file for verified peer lines")
	report := flag.String("report", "", "output file for JSON verification report")
	min := flag.Int("min", 8, "minimum verified peers required (else exit 1)")
	concurrency := flag.Int("concurrency", 16, "max concurrent dials")
	timeoutSec := flag.Int("timeout", 5, "per-peer dial+handshake timeout in seconds")
	flag.Parse()

	if *in == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "usage: verifypeer -in candidates.txt -out peers.txt [-report report.json] [-min N]")
		os.Exit(2)
	}

	candidates, err := readCandidates(*in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read candidates: %v\n", err)
		os.Exit(1)
	}

	results := verifyAll(candidates, *concurrency, time.Duration(*timeoutSec)*time.Second)
	verified := verifiedLines(results)
	sort.Strings(verified) // deterministic output independent of dial timing

	if err := os.WriteFile(*out, []byte(strings.Join(verified, "\n")+"\n"), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write out: %v\n", err)
		os.Exit(1)
	}

	if *report != "" {
		b, _ := json.MarshalIndent(results, "", "  ")
		_ = os.WriteFile(*report, b, 0o644)
	}

	mm := mismatches(results)
	fmt.Printf("candidates=%d verified=%d unreachable=%d id_mismatch=%d invalid=%d\n",
		len(candidates), len(verified),
		countStatus(results, StatusUnreachable), len(mm), countStatus(results, StatusInvalid))
	for _, r := range mm {
		fmt.Fprintf(os.Stderr, "ID MISMATCH: claimed=%s remote=%s line=%s\n", r.ClaimedID, r.RemoteID, r.Line)
	}

	if len(verified) < *min {
		fmt.Fprintf(os.Stderr, "FAIL: only %d verified peers, need >= %d\n", len(verified), *min)
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
