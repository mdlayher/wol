wol
===

Command `wol` is a simple Wake-on-LAN client.  It can issue Wake-on-LAN magic
packets using both UDP or raw ethernet sockets.

Usage
-----

```
Usage of ./wol:
  -a="": network address for Wake-on-LAN magic packet
  -i="": network interface to use to send Wake-on-LAN magic packet
  -p="": optional password for Wake-on-LAN magic packet
  -t="": target for Wake-on-LAN magic packet
```

Issue Wake-on-LAN magic packet using UDP network address:

```
$ ./wol -a 192.168.1.1:7 -t 00:12:7f:eb:6b:40
```

Issue Wake-on-LAN magic packet using raw ethernet sockets (requires elevated
privileges):

```
$ sudo ./wol -i eth0 -t 00:12:7f:eb:6b:40
```
