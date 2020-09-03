package wol

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClientWakePassword(t *testing.T) {
	var tests = []struct {
		name     string
		target   net.HardwareAddr
		password []byte
		err      error
	}{
		{
			name:   "5 byte target",
			target: make(net.HardwareAddr, 5),
			err:    errInvalidTarget,
		},
		{
			name:   "7 byte target",
			target: make(net.HardwareAddr, 7),
			err:    errInvalidTarget,
		},
		{
			name:     "5 bytes password",
			target:   make(net.HardwareAddr, 6),
			password: make([]byte, 5),
			err:      ErrInvalidPassword,
		},
		{
			name:     "7 byte password",
			target:   make(net.HardwareAddr, 6),
			password: make([]byte, 7),
			err:      ErrInvalidPassword,
		},
		{
			name:     "OK, no password",
			target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
			password: []byte{},
		},
		{
			name:     "OK, 4 byte password",
			target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
			password: []byte{1, 2, 3, 4},
		},
		{
			name:     "OK, 6 byte password",
			target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
			password: []byte{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeToPacketConn{}
			c := &Client{
				p: p,
			}

			// Address hardcoded because it doesn't matter for tests.
			err := c.WakePassword("127.0.0.1:0", tt.target, tt.password)

			if err != nil && tt.err == nil {
				t.Fatal("expected an error, but none occurred")
			}
			if err == nil && tt.err != nil {
				t.Fatalf("failed to send: %v", err)
			}
			if err != nil {
				return
			}

			mp := new(MagicPacket)
			if err := mp.UnmarshalBinary(p.b); err != nil {
				t.Fatalf("failed to unmarshal MagicPacket: %v", err)
			}

			wantMP := &MagicPacket{
				Target:   tt.target,
				Password: tt.password,
			}

			if diff := cmp.Diff(wantMP, mp); diff != "" {
				t.Fatalf("unexpected MagicPacket (-want +got):\n%s", diff)
			}
		})
	}
}
