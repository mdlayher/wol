// Package wol implements a Wake-on-LAN client.
package wol

import (
	"bytes"
	"errors"
	"io"
	"net"
)

const (
	// EtherType is the registered EtherType for Wake-on-LAN over Ethernet.
	// See: https://wiki.wireshark.org/WakeOnLAN.
	EtherType = 0x0842
)

var (
	// errInvalidPassword is returned if a MagicPacket's Password field is
	// not exactly 0 (empty), 4, or 6 bytes in length.
	errInvalidPassword = errors.New("invalid password length")

	// errInvalidSyncStream is returned if a MagicPacket's synchronization
	// stream is incorrect.
	errInvalidSyncStream = errors.New("invalid synchronization stream")

	// errInvalidTarget is returned if a MagicPacket does not contain the
	// same target hardware address repeated 16 times.
	errInvalidTarget = errors.New("invalid hardware address target")
)

var (
	// syncStream is a 6 byte slice which always occurs at the beginning of a
	// Wake-on-LAN magic packet.
	syncStream = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// A MagicPacket is a Wake-on-LAN packet.  It specifies a target hardware
// address to wake, and optionally, a password used to authenticate the
// MagicPacket wake request.
type MagicPacket struct {
	// Target specifies the hardware address of a LAN device to wake using
	// this MagicPacket.
	Target net.HardwareAddr

	// Password specifies an optional password for this MagicPacket.  Password
	// must be exactly 0 (empty), 4, or 6 bytes in length.
	Password []byte
}

// MarshalBinary allocates a byte slice and marshals a MagicPacket into binary
// form.
//
// If p.Target is not exactly 6 bytes in length, errInvalidTarget is returned.
//
// If p.Password is not exactly 0 (empty), 4, or 6 bytes in length,
// errInvalidPassword is returned.
func (p *MagicPacket) MarshalBinary() ([]byte, error) {
	// Must be 6 byte ethernet hardware address
	if len(p.Target) != 6 {
		return nil, errInvalidTarget
	}

	// Verify password is correct length
	if pl := len(p.Password); pl != 0 && pl != 4 && pl != 6 {
		return nil, errInvalidPassword
	}

	//    6 bytes: synchronization stream
	// 6*16 bytes: repeated target ethernet hardware address
	//    N bytes: password
	b := make([]byte, 6+(len(p.Target)*16)+len(p.Password))

	// Synchronization stream must always be present first
	copy(b[0:6], syncStream)

	// Place repeated target hardware address 16 times
	hl := len(p.Target)
	for i := 0; i < 16; i++ {
		copy(b[6+(hl*i):6+(hl*i)+hl], p.Target)
	}

	// Add password at end of slice
	copy(b[len(b)-len(p.Password):], p.Password)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a MagicPacket.
//
// If the byte slice does not contain enough data to unmarshal a valid
// MagicPacket, io.ErrUnexpectedEOF is returned.
func (p *MagicPacket) UnmarshalBinary(b []byte) error {
	// Must contain sync stream and 16 repeated targets
	if len(b) < 6+(6*16) {
		return io.ErrUnexpectedEOF
	}

	// Sync stream must be correct
	if !bytes.Equal(b[0:6], syncStream) {
		return errInvalidSyncStream
	}

	// Hardware address must correctly repeat 16 times
	for i := 0; i < 16; i++ {
		if !bytes.Equal(b[6:12], b[6+(6*i):6+(6*i)+6]) {
			return errInvalidTarget
		}
	}

	pl := len(b[6+(6*16):])

	// Allocate a single byte slice for target and password, and
	// reslice it to store fields
	bb := make([]byte, 6+pl)

	copy(bb[0:6], b[6:12])
	p.Target = bb[0:6]

	copy(bb[6:], b[len(b)-pl:])
	p.Password = bb[6:]

	// Password must be 0 (empty), 4, or 6 bytes in length
	if pl != 0 && pl != 4 && pl != 6 {
		return errInvalidPassword
	}

	return nil
}
