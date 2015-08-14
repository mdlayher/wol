package wol

import (
	"net"
)

// A Client is a Wake-on-LAN client which utilizes a UDP socket.  It can be
// used to send WoL magic packets to other machines using their network
// address.
type Client struct{}

// Wake sends a Wake-on-LAN magic packet to a device with the specified
// network and hardware address.
//
// If target is not a 6 byte Ethernet hardware address, ErrInvalidTarget
// is returned.
func (c *Client) Wake(addr string, target net.HardwareAddr) error {
	return c.WakePassword(addr, target, nil)
}

// WakePassword sends a Wake-on-LAN magic packet to a device with the
// specified network and hardware address, using the specified password.
//
// If target is not a 6 byte Ethernet hardware address, ErrInvalidTarget
// is returned.
//
// The password must be exactly 0 (empty), 4, or 6 bytes in length, or
// ErrInvalidPassword will be returned.
func (c *Client) WakePassword(addr string, target net.HardwareAddr, password []byte) error {
	return c.withConn(addr, func(p net.Conn) error {
		return c.sendWake(p, target, password)
	})
}

// sendWake crafts a magic packet using the input parameters and sends the
// frame over a UDP socket to attempt to wake a machine.
func (c *Client) sendWake(p net.Conn, target net.HardwareAddr, password []byte) error {
	// Create magic packet with target and password
	mp := &MagicPacket{
		Target:   target,
		Password: password,
	}
	mpb, err := mp.MarshalBinary()
	if err != nil {
		return err
	}

	// Send magic packet to target over UDP socket
	_, err = p.Write(mpb)
	return err
}

// withConn resolves address addr, opens a UDP socket, and passes the socket
// as a parameter to the input closure.  The socket is closed once the closure
// returns.
func (c *Client) withConn(addr string, fn func(p net.Conn) error) error {
	// Resolve destination address
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// Dial a UDP connection
	p, err := net.DialUDP("udp", nil, uaddr)
	if err != nil {
		return err
	}
	defer p.Close()

	// Invoke closure with connection
	return fn(p)
}
