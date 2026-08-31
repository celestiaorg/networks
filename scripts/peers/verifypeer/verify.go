package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/libs/protoio"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
	tmp2p "github.com/cometbft/cometbft/proto/tendermint/p2p"
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
	StatusVerified        Status = "verified"
	StatusIDMismatch      Status = "id_mismatch"
	StatusWrongChain      Status = "wrong_chain"
	StatusHandshakeFailed Status = "handshake_failed"
	StatusUnreachable     Status = "unreachable"
	StatusInvalid         Status = "invalid"
)

type Result struct {
	Line      string `json:"line"`
	ClaimedID string `json:"claimed_id"`
	RemoteID  string `json:"remote_id,omitempty"`
	Network   string `json:"network,omitempty"`
	Status    Status `json:"status"`
	Err       string `json:"error,omitempty"`
}

// VerifyPeer dials the peer and runs the full CometBFT peer handshake: a secret
// connection (which authenticates the remote node ID against the claimed id@)
// followed by a NodeInfo exchange (which reveals the chain the peer is actually
// serving). A peer is only verified when both the node ID and the chain match.
//
// The statuses are deliberately distinct. "unreachable" means the TCP dial
// failed; "handshake_failed" means the port is open but did not complete the
// CometBFT handshake. Collapsing the two hides whether a refresh found dead
// hosts or live hosts that rejected us.
func VerifyPeer(line string, chainID string, timeout time.Duration) Result {
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
	// Deadline covers the rest of the handshake; total time may be up to ~2×timeout.
	_ = c.SetDeadline(time.Now().Add(timeout))

	// Use a one-time key; we only need to authenticate the remote peer, not ourselves.
	ourKey := ed25519.GenPrivKey()
	sc, err := conn.MakeSecretConnection(c, ourKey)
	if err != nil {
		res.Status = StatusHandshakeFailed
		res.Err = "secret connection: " + err.Error()
		return res
	}
	remoteID := p2p.PubKeyToID(sc.RemotePubKey())
	res.RemoteID = string(remoteID)
	if !strings.EqualFold(string(remoteID), string(claimedID)) {
		res.Status = StatusIDMismatch
		return res
	}

	nodeInfo, err := exchangeNodeInfo(sc, ourKey, chainID)
	if err != nil {
		res.Status = StatusHandshakeFailed
		res.Err = "node info: " + err.Error()
		return res
	}
	res.Network = nodeInfo.Network
	if chainID != "" && nodeInfo.Network != chainID {
		res.Status = StatusWrongChain
		return res
	}
	res.Status = StatusVerified
	return res
}

// exchangeNodeInfo performs the NodeInfo half of the CometBFT peer handshake and
// returns the remote node's NodeInfo. Mirrors the unexported handshake() in
// cometbft/p2p/transport.go: both sides write and read concurrently, so neither
// blocks waiting for the other to speak first.
func exchangeNodeInfo(sc *conn.SecretConnection, ourKey ed25519.PrivKey, chainID string) (p2p.DefaultNodeInfo, error) {
	// Advertise a plausible, compatible NodeInfo so the peer has no reason to
	// hang up before sending its own.
	ours := p2p.DefaultNodeInfo{
		ProtocolVersion: p2p.ProtocolVersion{P2P: 8, Block: 11, App: 0},
		DefaultNodeID:   p2p.PubKeyToID(ourKey.PubKey()),
		ListenAddr:      "0.0.0.0:26656",
		Network:         chainID,
		Version:         "0.40.2",
		Channels:        []byte{0x40},
		Moniker:         "networks-peer-verifier",
	}

	errc := make(chan error, 2)
	var pb tmp2p.DefaultNodeInfo
	go func() {
		_, err := protoio.NewDelimitedWriter(sc).WriteMsg(ours.ToProto())
		errc <- err
	}()
	go func() {
		_, err := protoio.NewDelimitedReader(sc, p2p.MaxNodeInfoSize()).ReadMsg(&pb)
		errc <- err
	}()
	for i := 0; i < cap(errc); i++ {
		if err := <-errc; err != nil {
			return p2p.DefaultNodeInfo{}, err
		}
	}
	return p2p.DefaultNodeInfoFromToProto(&pb)
}
