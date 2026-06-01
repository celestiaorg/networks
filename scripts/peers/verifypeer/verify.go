package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
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

type Status string

const (
	StatusVerified    Status = "verified"
	StatusIDMismatch  Status = "id_mismatch"
	StatusUnreachable Status = "unreachable"
	StatusInvalid     Status = "invalid"
)

type Result struct {
	Line      string `json:"line"`
	ClaimedID string `json:"claimed_id"`
	RemoteID  string `json:"remote_id,omitempty"`
	Status    Status `json:"status"`
	Err       string `json:"error,omitempty"`
}

// VerifyPeer dials the peer, performs the CometBFT secret-connection handshake,
// and confirms the authenticated remote node ID equals the claimed id@.
func VerifyPeer(line string, timeout time.Duration) Result {
	claimedID, dialAddr, err := ParsePeer(line)
	if err != nil {
		return Result{Line: line, Status: StatusInvalid, Err: err.Error()}
	}
	res := Result{Line: line, ClaimedID: string(claimedID)}

	d := net.Dialer{Timeout: timeout}
	c, err := d.Dial("tcp", dialAddr)
	if err != nil {
		res.Status = StatusUnreachable
		res.Err = err.Error()
		return res
	}
	defer c.Close()
	// Deadline covers the handshake; total time may be up to ~2×timeout.
	_ = c.SetDeadline(time.Now().Add(timeout))

	// Use a one-time key; we only need to authenticate the remote peer, not ourselves.
	sc, err := conn.MakeSecretConnection(c, ed25519.GenPrivKey())
	if err != nil {
		res.Status = StatusUnreachable
		res.Err = "handshake: " + err.Error()
		return res
	}
	remoteID := p2p.PubKeyToID(sc.RemotePubKey())
	res.RemoteID = string(remoteID)
	if !strings.EqualFold(string(remoteID), string(claimedID)) {
		res.Status = StatusIDMismatch
		return res
	}
	res.Status = StatusVerified
	return res
}
