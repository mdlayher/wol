# wol [![builds.sr.ht status](https://builds.sr.ht/~mdlayher/wol.svg)](https://builds.sr.ht/~mdlayher/wol?) [![GoDoc](https://godoc.org/github.com/mdlayher/wol?status.svg)](https://godoc.org/github.com/mdlayher/wol) [![Go Report Card](https://goreportcard.com/badge/github.com/mdlayher/wol)](https://goreportcard.com/report/github.com/mdlayher/wol)

Package `wol` implements a Wake-on-LAN client. MIT Licensed.

This package exposes two types, which operate slightly differently:

- `Client`: WoL client which uses UDP sockets to send magic packets
- `RawClient` WoL client which uses raw Ethernet sockets to send magic packets

For most use cases, the `Client` type will be sufficient.  The `RawClient` type
requires elevated privileges (root user) and works on Linux or *BSD/macOS only.
