go-scp [![CI](https://github.com/mondoohq/go-scp/actions/workflows/ci.yml/badge.svg)](https://github.com/mondoohq/go-scp/actions/workflows/ci.yml)  [![Go Report Card](https://goreportcard.com/badge/github.com/hnakamur/go-scp)](https://goreportcard.com/report/github.com/hnakamur/go-scp) [![PkgGoDev](https://pkg.go.dev/badge/github.com/hnakamur/go-scp)](https://pkg.go.dev/github.com/hnakamur/go-scp) [![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)
======
## About this fork

This is Mondoo's fork of [hnakamur/go-scp](https://github.com/hnakamur/go-scp), which the
original author stopped maintaining. It is maintained here, and kept building and tested
across the platforms Go supports.

OpenSSH has [deprecated the scp protocol](https://lwn.net/Articles/835962/) in favour of
SFTP, so prefer SFTP for new work. The scp path remains useful for reaching hosts where
the SFTP subsystem is unavailable, which is why this library is still here.

## Usage
A scp client library written in Go.
The remote server must have the scp command.

## Example
Please refer to [the example at godoc](https://godoc.org/github.com/hnakamur/go-scp#example-package).

## License
MIT
