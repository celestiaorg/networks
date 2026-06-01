package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cometbft/cometbft/p2p"
)

// ParsePeer splits an "id@host:port" line into a lowercased node ID and a
// dial address. The ID must be 40 hex characters (a 20-byte CometBFT node ID).
// The returned dial address is not validated here; validation happens at dial time.
func ParsePeer(line string) (id p2p.ID, dialAddr string, err error) {
	line = strings.TrimSpace(line)
	at := strings.IndexByte(line, '@')
	if at < 0 {
		return "", "", fmt.Errorf("missing '@' in peer %q", line)
	}
	rawID := strings.ToLower(line[:at])
	addr := line[at+1:]
	if rawID == "" || addr == "" {
		return "", "", fmt.Errorf("empty id or address in peer %q", line)
	}
	if len(rawID) != 40 {
		return "", "", fmt.Errorf("id must be 40 hex chars, got %d in peer %q", len(rawID), line)
	}
	if _, err := hex.DecodeString(rawID); err != nil {
		return "", "", fmt.Errorf("id must be hex in peer %q: %w", line, err)
	}
	return p2p.ID(rawID), addr, nil
}
