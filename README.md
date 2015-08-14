wol [![Build Status](https://travis-ci.org/mdlayher/wol.svg?branch=master)](https://travis-ci.org/mdlayher/wol) [![Coverage Status](https://coveralls.io/repos/mdlayher/wol/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/wol?branch=master) [![GoDoc](http://godoc.org/github.com/mdlayher/wol?status.svg)](http://godoc.org/github.com/mdlayher/wol)
===

Package `wol` implements a Wake-on-LAN client.  MIT Licensed.

This package exposes two types, which operate slightly differently:
- `Client`: WoL client which uses UDP sockets to send magic packets
- `RawClient` WoL client which uses raw Ethernet sockets to send magic packets

For most use cases, the `Client` type will be sufficient.  The `RawClient` type requires
elevated privileges (root user) and currently works with Linux only.
