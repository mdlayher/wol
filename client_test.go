package wol

import (
	"bytes"
	"net"
	"testing"
)

func TestClientWakePassword(t *testing.T) {
	var tests = []struct {
		desc     string
		target   net.HardwareAddr
		password []byte
		err      error
	}{
		{
			desc:   "5 byte target",
			target: make(net.HardwareAddr, 5),
			err:    ErrInvalidTarget,
		},
		{
			desc:   "7 byte target",
			target: make(net.HardwareAddr, 7),
			err:    ErrInvalidTarget,
		},
		{
			desc:     "5 bytes password",
			target:   make(net.HardwareAddr, 6),
			password: make([]byte, 5),
			err:      ErrInvalidPassword,
		},
		{
			desc:     "7 byte password",
			target:   make(net.HardwareAddr, 6),
			password: make([]byte, 7),
			err:      ErrInvalidPassword,
		},
		{
			desc:   "OK, no password",
			target: net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
		},
		{
			desc:     "OK, 4 byte password",
			target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
			password: []byte{1, 2, 3, 4},
		},
		{
			desc:     "OK, 6 byte password",
			target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
			password: []byte{1, 2, 3, 4, 5, 6},
		},
	}

	for i, tt := range tests {
		p := &writeToPacketConn{}
		c := &Client{
			p: p,
		}

		// Address hardcoded because it doesn't matter for tests
		if err := c.WakePassword("127.0.0.1:0", tt.target, tt.password); err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		mp := new(MagicPacket)
		if err := mp.UnmarshalBinary(p.b); err != nil {
			t.Fatal(err)
		}

		if want, got := tt.target, mp.Target; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected MagicPacket target:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
		if want, got := len(tt.password), len(mp.Password); want != got {
			t.Fatalf("[%02d] test %q, unexpected MagicPacket password length:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}
