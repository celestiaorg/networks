package main

import (
	"net"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
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

// startHandshakeServer accepts one connection, performs the server side of the
// secret-connection handshake using serverKey, then closes. Returns the
// listener address and the server's node ID.
func startHandshakeServer(t *testing.T, serverKey ed25519.PrivKey) (addr string, nodeID p2p.ID) {
	t.Helper()
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
		defer c.Close()
		_ = c.SetDeadline(time.Now().Add(3 * time.Second))
		// Server side of the handshake. Result is discarded; we only need the
		// dialer to complete its side.
		_, _ = conn.MakeSecretConnection(c, serverKey)
	}()
	return ln.Addr().String(), p2p.PubKeyToID(serverKey.PubKey())
}

func TestVerifyPeer_Match(t *testing.T) {
	key := ed25519.GenPrivKey()
	addr, id := startHandshakeServer(t, key)
	res := VerifyPeer(string(id)+"@"+addr, 3*time.Second)
	if res.Status != StatusVerified {
		t.Fatalf("status = %q (%s), want verified", res.Status, res.Err)
	}
	if res.RemoteID != string(id) {
		t.Fatalf("remoteID = %q, want %q", res.RemoteID, id)
	}
}

func TestVerifyPeer_Mismatch(t *testing.T) {
	key := ed25519.GenPrivKey()
	addr, _ := startHandshakeServer(t, key)
	wrongID := "0000000000000000000000000000000000000000"
	res := VerifyPeer(wrongID+"@"+addr, 3*time.Second)
	if res.Status != StatusIDMismatch {
		t.Fatalf("status = %q, want id_mismatch", res.Status)
	}
}

func TestVerifyPeer_Unreachable(t *testing.T) {
	// Reserve a port then close it so the dial is refused/times out.
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	_ = ln.Close()
	id := "1111111111111111111111111111111111111111"
	res := VerifyPeer(id+"@"+addr, 1*time.Second)
	if res.Status != StatusUnreachable {
		t.Fatalf("status = %q, want unreachable", res.Status)
	}
}
