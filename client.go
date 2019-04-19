package wol

import (
	"net"
)

// A Client is a Wake-on-LAN client which utilizes a UDP socket.  It can be
// used to send WoL magic packets to other machines using their network
// address.
type Client struct {
	p net.PacketConn
}

// NewClient creates a new Client which binds to any available UDP port to
// send Wake-on-LAN magic packets.
func NewClient() (*Client, error) {
	// Bind to any available UDP port.
	p, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, err
	}

	return &Client{
		p: p,
	}, nil
}

// Close closes a Client's UDP socket.
func (c *Client) Close() error {
	return c.p.Close()
}

// Wake sends a Wake-on-LAN magic packet to a device with the specified
// network and hardware address.
func (c *Client) Wake(addr string, target net.HardwareAddr) error {
	return c.WakePassword(addr, target, nil)
}

// WakePassword sends a Wake-on-LAN magic packet to a device with the
// specified network and hardware address, using the specified password.
//
// The password must be exactly 0 (empty), 4, or 6 bytes in length.
func (c *Client) WakePassword(addr string, target net.HardwareAddr, password []byte) error {
	return c.sendWake(addr, target, password)
}

// sendWake crafts a magic packet using the input parameters and sends the
// packet over a UDP socket to attempt to wake a machine.
func (c *Client) sendWake(addr string, target net.HardwareAddr, password []byte) error {
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// Create magic packet with target and password.
	mp := &MagicPacket{
		Target:   target,
		Password: password,
	}
	mpb, err := mp.MarshalBinary()
	if err != nil {
		return err
	}

	// Send magic packet to target over UDP socket.
	_, err = c.p.WriteTo(mpb, uaddr)
	return err
}
