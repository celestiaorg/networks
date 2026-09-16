package main

import (
	"net"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/libs/protoio"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
	tmp2p "github.com/cometbft/cometbft/proto/tendermint/p2p"
)

func TestParsePeer(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantID   string
		wantAddr string
		wantErr  bool
	}{
		{"valid", "13612522d3ce71a181881370b8f40103d4def9f1@84.32.32.148:16400", "13612522d3ce71a181881370b8f40103d4def9f1", "84.32.32.148:16400", false},
		{"uppercase id lowercased", "13612522D3CE71A181881370B8F40103D4DEF9F1@1.2.3.4:26656", "13612522d3ce71a181881370b8f40103d4def9f1", "1.2.3.4:26656", false},
		{"missing @", "13612522d3ce71a181881370b8f40103d4def9f1-1.2.3.4:26656", "", "", true},
		{"empty id", "@1.2.3.4:26656", "", "", true},
		{"empty addr", "13612522d3ce71a181881370b8f40103d4def9f1@", "", "", true},
		{"short id", "abc@1.2.3.4:26656", "", "", true},
		{"leading/trailing whitespace trimmed", "  13612522d3ce71a181881370b8f40103d4def9f1@1.2.3.4:26656\n", "13612522d3ce71a181881370b8f40103d4def9f1", "1.2.3.4:26656", false},
		{"long id (41 hex chars)", "13612522d3ce71a181881370b8f40103d4def9f1a@1.2.3.4:26656", "", "", true},
		{"non-hex id (40 chars with z)", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz@1.2.3.4:26656", "", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, addr, err := ParsePeer(tc.line)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if string(id) != tc.wantID || addr != tc.wantAddr {
				t.Fatalf("got (%q, %q), want (%q, %q)", id, addr, tc.wantID, tc.wantAddr)
			}
		})
	}
}

// startPeerServer accepts one connection and plays the server side of the full
// CometBFT peer handshake: secret connection with serverKey, then a NodeInfo
// exchange advertising network. Returns the listener address and node ID.
func startPeerServer(t *testing.T, serverKey ed25519.PrivKey, network string) (addr string, nodeID p2p.ID) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	id := p2p.PubKeyToID(serverKey.PubKey())
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		_ = c.SetDeadline(time.Now().Add(5 * time.Second))
		sc, err := conn.MakeSecretConnection(c, serverKey)
		if err != nil {
			return
		}
		ni := p2p.DefaultNodeInfo{
			ProtocolVersion: p2p.ProtocolVersion{P2P: 8, Block: 11, App: 1},
			DefaultNodeID:   id,
			ListenAddr:      "0.0.0.0:26656",
			Network:         network,
			Version:         "0.40.2",
			Channels:        []byte{0x40},
			Moniker:         "test-server",
		}
		_, _ = protoio.NewDelimitedWriter(sc).WriteMsg(ni.ToProto())
		// Drain the client's NodeInfo so it does not block on a full buffer.
		var peer tmp2p.DefaultNodeInfo
		_, _ = protoio.NewDelimitedReader(sc, p2p.MaxNodeInfoSize()).ReadMsg(&peer)
	}()
	return ln.Addr().String(), id
}

func TestVerifyPeer_Match(t *testing.T) {
	key := ed25519.GenPrivKey()
	addr, id := startPeerServer(t, key, "celestia")
	res := VerifyPeer(string(id)+"@"+addr, "celestia", 3*time.Second)
	if res.Status != StatusVerified {
		t.Fatalf("status = %q (%s), want verified", res.Status, res.Err)
	}
	if res.RemoteID != string(id) {
		t.Fatalf("remoteID = %q, want %q", res.RemoteID, id)
	}
	if res.Network != "celestia" {
		t.Fatalf("network = %q, want %q", res.Network, "celestia")
	}
}

func TestVerifyPeer_WrongChain(t *testing.T) {
	key := ed25519.GenPrivKey()
	// Node ID matches, but the peer is on mocha-5 while we asked for celestia.
	addr, id := startPeerServer(t, key, "mocha-5")
	res := VerifyPeer(string(id)+"@"+addr, "celestia", 3*time.Second)
	if res.Status != StatusWrongChain {
		t.Fatalf("status = %q (%s), want wrong_chain", res.Status, res.Err)
	}
	if res.Network != "mocha-5" {
		t.Fatalf("network = %q, want %q", res.Network, "mocha-5")
	}
}

func TestVerifyPeer_Mismatch(t *testing.T) {
	key := ed25519.GenPrivKey()
	addr, _ := startPeerServer(t, key, "celestia")
	wrongID := "0000000000000000000000000000000000000000"
	res := VerifyPeer(wrongID+"@"+addr, "celestia", 3*time.Second)
	if res.Status != StatusIDMismatch {
		t.Fatalf("status = %q, want id_mismatch", res.Status)
	}
}

func TestVerifyPeer_Unreachable(t *testing.T) {
	// Reserve a loopback port then close it so the dial is refused. The brief
	// reuse window is acceptable for this hermetic test.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	id := "1111111111111111111111111111111111111111"
	res := VerifyPeer(id+"@"+addr, "celestia", 1*time.Second)
	if res.Status != StatusUnreachable {
		t.Fatalf("status = %q, want unreachable", res.Status)
	}
}

// A TCP port that accepts but speaks no CometBFT is a handshake failure, not an
// unreachable host: the distinction matters when triaging why a refresh failed.
func TestVerifyPeer_HandshakeFailedIsNotUnreachable(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		// Accept, say nothing, hang up: TCP is open, the handshake cannot complete.
		_ = c.Close()
	}()
	id := "1111111111111111111111111111111111111111"
	res := VerifyPeer(id+"@"+ln.Addr().String(), "celestia", 2*time.Second)
	if res.Status != StatusHandshakeFailed {
		t.Fatalf("status = %q, want handshake_failed", res.Status)
	}
}
