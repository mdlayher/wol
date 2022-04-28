package wol

import (
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/ethernet"
)

func TestRawClientWakePassword(t *testing.T) {
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

	zeroHW := make(net.HardwareAddr, 6)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeToPacketConn{}
			c := &RawClient{
				ifi: &net.Interface{
					HardwareAddr: zeroHW,
				},
				p: p,
			}

			err := c.WakePassword(tt.target, tt.password)

			if err != nil && tt.err == nil {
				t.Fatal("expected an error, but none occurred")
			}
			if err == nil && tt.err != nil {
				t.Fatalf("failed to send: %v", err)
			}
			if err != nil {
				return
			}

			f := new(ethernet.Frame)
			if err := f.UnmarshalBinary(p.b); err != nil {
				t.Fatalf("failed to unmarshal Ethernet frame: %v", err)
			}

			wantEth := &ethernet.Frame{
				Destination: tt.target,
				Source:      zeroHW,
				EtherType:   EtherType,
			}

			// Copy out payload and nil out for base Ethernet frame comparison.
			pl := make([]byte, len(f.Payload))
			copy(pl, f.Payload)
			f.Payload = nil

			if diff := cmp.Diff(wantEth, f); diff != "" {
				t.Fatalf("unexpected Ethernet frame (-want +got):\n%s", diff)
			}

			mp := new(MagicPacket)
			if err := mp.UnmarshalBinary(pl); err != nil {
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
