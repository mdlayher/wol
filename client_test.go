package wol

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestClient_sendWake(t *testing.T) {
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
		p := &writeConn{}
		c := &Client{}

		if err := c.sendWake(p, tt.target, tt.password); err != nil || tt.err != nil {
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

type writeConn struct {
	b []byte
	noopConn
}

func (w *writeConn) Write(b []byte) (int, error) {
	w.b = make([]byte, len(b))
	copy(w.b, b)

	return len(b), nil
}

// noopConn is a net.Conn which simply no-ops any input.  It is
// embedded in other implementations so they do not have to implement every
// single method.
type noopConn struct{}

func (noopConn) Read(b []byte) (int, error)  { return 0, nil }
func (noopConn) Write(b []byte) (int, error) { return 0, nil }

func (noopConn) Close() error                       { return nil }
func (noopConn) LocalAddr() net.Addr                { return nil }
func (noopConn) RemoteAddr() net.Addr               { return nil }
func (noopConn) SetDeadline(t time.Time) error      { return nil }
func (noopConn) SetReadDeadline(t time.Time) error  { return nil }
func (noopConn) SetWriteDeadline(t time.Time) error { return nil }
