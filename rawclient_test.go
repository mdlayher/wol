package wol

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/mdlayher/ethernet"
)

func TestRawClientWakePassword(t *testing.T) {
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

	zeroHW := make(net.HardwareAddr, 6)

	for i, tt := range tests {
		p := &writeToPacketConn{}
		c := &RawClient{
			ifi: &net.Interface{
				HardwareAddr: zeroHW,
			},
			p: p,
		}

		if err := c.WakePassword(tt.target, tt.password); err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		// Special case, trim trailing four zero characters, which would be occupied
		// by an Ethernet frame check sequence; if one existed
		if bytes.Equal(p.b[len(p.b)-4:], make([]byte, 4)) {
			p.b = p.b[:len(p.b)-4]
		}

		f := new(ethernet.Frame)
		if err := f.UnmarshalBinary(p.b); err != nil {
			t.Fatal(err)
		}

		if want, got := tt.target, f.Destination; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Ethernet frame destination:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
		if want, got := zeroHW, f.Source; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Ethernet frame source:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
		if want, got := EtherType, f.EtherType; want != got {
			t.Fatalf("[%02d] test %q, unexpected Ethernet frame EtherType:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}

		mp := new(MagicPacket)
		if err := mp.UnmarshalBinary(f.Payload); err != nil {
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

type writeToPacketConn struct {
	b []byte
	noopPacketConn
}

func (w *writeToPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	w.b = make([]byte, len(b))
	copy(w.b, b)

	return len(b), nil
}

// noopPacketConn is a net.PacketConn which simply no-ops any input.  It is
// embedded in other implementations so they do not have to implement every
// single method.
type noopPacketConn struct{}

func (noopPacketConn) ReadFrom(b []byte) (int, net.Addr, error)     { return 0, nil, nil }
func (noopPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) { return 0, nil }

func (noopPacketConn) Close() error                       { return nil }
func (noopPacketConn) LocalAddr() net.Addr                { return nil }
func (noopPacketConn) SetDeadline(t time.Time) error      { return nil }
func (noopPacketConn) SetReadDeadline(t time.Time) error  { return nil }
func (noopPacketConn) SetWriteDeadline(t time.Time) error { return nil }
