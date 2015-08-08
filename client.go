package wol

import (
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

// A Client is a Wake-on-LAN client.  It can be used to send WoL magic packets
// to other machines on a local network, using their hardware addresses.
type Client struct {
	ifi *net.Interface
	p   net.PacketConn
}

// NewClient creates a new Client using the specified network interface.
func NewClient(iface string) (*Client, error) {
	// Verify interface exists
	ifi, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	// Open raw socket to send Wake-on-LAN magic packets
	p, err := raw.ListenPacket(ifi, raw.ProtocolWoL)
	if err != nil {
		return nil, err
	}

	return &Client{
		ifi: ifi,
		p:   p,
	}, nil
}

// Close closes a Client's raw socket.
func (c *Client) Close() error {
	return c.p.Close()
}

// Wake sends a Wake-on-LAN magic packet to the specified hardware address.
//
// If target is not a 6 byte Ethernet hardware address, ErrInvalidTarget
// is returned.
func (c *Client) Wake(target net.HardwareAddr) error {
	return c.sendWake(target, nil)
}

// WakePassword sends a Wake-on-LAN magic packet to the specified hardware
// address, using the specified Password.
//
// If target is not a 6 byte Ethernet hardware address, ErrInvalidTarget
// is returned.
//
// The password must be exactly 0 (empty), 4, or 6 bytes in length, or
// ErrInvalidPassword will be returned.
func (c *Client) WakePassword(target net.HardwareAddr, password []byte) error {
	return c.sendWake(target, password)
}

// sendWake crafts a magic packet using the input parameters, stores it in an
// Ethernet frame, and sends the frame over a raw socket to attempt to wake
// a machine.
func (c *Client) sendWake(target net.HardwareAddr, password []byte) error {
	// Create magic packet with target and password
	p := &MagicPacket{
		Target:   target,
		Password: password,
	}
	pb, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	// Create Ethernet frame to carry magic packet
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

	// Send magic packet to target
	_, err = c.p.WriteTo(fb, &raw.Addr{
		HardwareAddr: target,
	})
	return err
}
