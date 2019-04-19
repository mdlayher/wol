# wol

Command `wol` is a simple Wake-on-LAN client.  It can issue Wake-on-LAN magic
packets using both UDP or Ethernet sockets.

## Usage

```text
$ ./wol -h
Usage of ./wol:
  -a string
        network address for Wake-on-LAN magic packet
  -i string
        network interface to use to send Wake-on-LAN magic packet
  -p string
        optional password for Wake-on-LAN magic packet
  -t string
        target for Wake-on-LAN magic packet
```

Issue Wake-on-LAN magic packet using UDP network address:

```text
./wol -a 192.168.1.1:7 -t 00:12:7f:eb:6b:40
```

Issue Wake-on-LAN magic packet using Ethernet sockets (requires elevated
privileges):

```text
sudo ./wol -i eth0 -t 00:12:7f:eb:6b:40
```
