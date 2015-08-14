// Command wol is a simple Wake-on-LAN client.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/mdlayher/wol"
)

var (
	addrFlag     = flag.String("a", "", "network address for Wake-on-LAN magic packet")
	ifaceFlag    = flag.String("i", "", "network interface to use to send Wake-on-LAN magic packet")
	targetFlag   = flag.String("t", "", "target for Wake-on-LAN magic packet")
	passwordFlag = flag.String("p", "", "optional password for Wake-on-LAN magic packet")
)

func main() {
	flag.Parse()

	// Validate hardware address
	target, err := net.ParseMAC(*targetFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Set password if one is present
	var password []byte
	if *passwordFlag != "" {
		password = []byte(*passwordFlag)
	}

	// Can only do raw or UDP mode, not both
	if *addrFlag != "" && *ifaceFlag != "" {
		log.Fatalf("must set '-a' or '-i' flag exclusively")
	}

	// Check for raw mode
	if *ifaceFlag != "" {
		if err := wakeRaw(*ifaceFlag, target, password); err != nil {
			log.Fatal(err)
		}

		log.Printf("sent raw Wake-on-LAN magic packet using %s to %s", *ifaceFlag, *targetFlag)
		return
	}

	// Use UDP mode
	if err := (&wol.Client{}).WakePassword(*addrFlag, target, password); err != nil {
		log.Fatal(err)
	}

	log.Printf("sent UDP Wake-on-LAN magic packet using %s to %s", *addrFlag, *targetFlag)
}

func wakeRaw(iface string, target net.HardwareAddr, password []byte) error {
	// Validate interface
	ifi, err := net.InterfaceByName(*ifaceFlag)
	if err != nil {
		return err
	}

	// Create client bound to specified interface
	c, err := wol.NewRawClient(ifi)
	if err != nil {
		return err
	}

	// Attempt to wake target machine
	if err := c.WakePassword(target, password); err != nil {
		log.Fatal(err)
	}

	return c.Close()
}
