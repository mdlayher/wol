package wol

import (
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/packet"
)

// A RawClient is a Wake-on-LAN client which operates directly on top of
// Ethernet frames using Ethernet sockets.  It can be used to send WoL magic
// packets to other machines on a local network, using their hardware addresses.
type RawClient struct {
	ifi *net.Interface
	p   net.PacketConn
}

// NewRawClient creates a new RawClient using the specified network interface.
//
// Note that Ethernet sockets typically require elevated user privileges, such
// as the 'root' user on Linux, or the 'SET_CAP_RAW' capability.
//
// For this reason, it is typically recommended to use the regular Client type
// instead, which operates over UDP.
func NewRawClient(ifi *net.Interface) (*RawClient, error) {
	// Open a packet socket to send Wake-on-LAN magic packets.
	// EtherType is set according to: https://wiki.wireshark.org/WakeOnLAN.
	p, err := packet.Listen(ifi, packet.Raw, EtherType, nil)
	if err != nil {
		return nil, err
	}

	return &RawClient{
		ifi: ifi,
		p:   p,
	}, nil
}

// Close closes a RawClient's socket.
func (c *RawClient) Close() error {
	return c.p.Close()
}

// Wake sends a Wake-on-LAN magic packet to the specified hardware address.
func (c *RawClient) Wake(target net.HardwareAddr) error {
	return c.WakePassword(target, nil)
}

// WakePassword sends a Wake-on-LAN magic packet to the specified hardware
// address, using the specified Password.
//
// The password must be exactly 0 (empty), 4, or 6 bytes in length.
func (c *RawClient) WakePassword(target net.HardwareAddr, password []byte) error {
	return c.sendWake(target, password)
}

// sendWake crafts a magic packet using the input parameters, stores it in an
// Ethernet frame, and sends the frame over an Ethernet socket to attempt to wake
// a machine.
func (c *RawClient) sendWake(target net.HardwareAddr, password []byte) error {
	// Create magic packet with target and password.
	p := &MagicPacket{
		Target:   target,
		Password: password,
	}
	pb, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	// Create Ethernet frame to carry magic packet.
	f := &ethernet.Frame{
		Destination: target,
		Source:      c.ifi.HardwareAddr,
		EtherType:   EtherType,
		Payload:     pb,
	}
	fb, err := f.MarshalBinary()
	if err != nil {
		return err
	}

	// Send magic packet to target.
	_, err = c.p.WriteTo(fb, &packet.Addr{
		HardwareAddr: target,
	})
	return err
}
