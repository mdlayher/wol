// Command wol is a simple Wake-on-LAN client.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/mdlayher/wol"
)

var (
	ifaceFlag    = flag.String("i", "eth0", "network interface to use to send Wake-on-LAN magic packet")
	targetFlag   = flag.String("t", "", "target for Wake-on-LAN magic packet")
	passwordFlag = flag.String("p", "", "optional password for Wake-on-LAN magic packet")
)

func main() {
	flag.Parse()

	// Validate interface
	ifi, err := net.InterfaceByName(*ifaceFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Create client bound to specified interface
	c, err := wol.NewRawClient(ifi)
	if err != nil {
		log.Fatal(err)
	}

	// Validate hardware address
	addr, err := net.ParseMAC(*targetFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Set password if one is present
	var password []byte
	if *passwordFlag != "" {
		password = []byte(*passwordFlag)
	}

	// Attempt to wake target machine
	if err := c.WakePassword(addr, password); err != nil {
		log.Fatal(err)
	}

	log.Printf("sent Wake-on-LAN magic packet using %s to %s", *ifaceFlag, *targetFlag)

	_ = c.Close()
}
