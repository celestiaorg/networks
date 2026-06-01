package main

import "testing"

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
