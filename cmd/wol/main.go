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

	target, err := net.ParseMAC(*targetFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Set password if one is present.
	var password []byte
	if *passwordFlag != "" {
		password = []byte(*passwordFlag)
	}

	// Can only do raw or UDP mode, not both.
	if *addrFlag != "" && *ifaceFlag != "" {
		log.Fatalf("must set '-a' or '-i' flag exclusively")
	}

	if *ifaceFlag != "" {
		if err := wakeRaw(*ifaceFlag, target, password); err != nil {
			log.Fatal(err)
		}

		log.Printf("sent raw Wake-on-LAN magic packet using %s to %s", *ifaceFlag, *targetFlag)
		return
	}

	if err := wakeUDP(*addrFlag, target, password); err != nil {
		log.Fatal(err)
	}

	log.Printf("sent UDP Wake-on-LAN magic packet using %s to %s", *addrFlag, *targetFlag)
}

func wakeRaw(iface string, target net.HardwareAddr, password []byte) error {
	ifi, err := net.InterfaceByName(*ifaceFlag)
	if err != nil {
		return err
	}

	c, err := wol.NewRawClient(ifi)
	if err != nil {
		return err
	}
	defer c.Close()

	// Attempt to wake target machine.
	return c.WakePassword(target, password)
}

func wakeUDP(addr string, target net.HardwareAddr, password []byte) error {
	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	// Attempt to wake target machine.
	return c.WakePassword(addr, target, password)
}
