# dht

[![CircleCI](https://circleci.com/gh/anacrolix/dht.svg?style=shield)](https://circleci.com/gh/anacrolix/dht)
[![Go Reference](https://pkg.go.dev/badge/github.com/anacrolix/dht/v2.svg)](https://pkg.go.dev/github.com/anacrolix/dht/v2)

This is a fork of [anacrolix dht](https://github.com/anacrolix/dht) adapting the library for use in [SCION Networks](https://www.scion-architecture.net/).
It was originally forked for usage in the [BitTorrent over SCION](https://github.com/martin31821/torrent) library.

## Installation

Get the library package with `go get github.com/anacrolix/dht/v2`, or the provided cmds with `go install github.com/anacrolix/dht/v2/cmd/...@latest`.

## Commands

Here I'll describe what some of the provided commands in `./cmd` do.

### dht-ping

Pings DHT nodes with the given network addresses.

    $ go run ./cmd/dht-ping router.bittorrent.com:6881 router.utorrent.com:6881
    2015/04/01 17:21:23 main.go:33: dht server on [::]:60058
    32f54e697351ff4aec29cdbaabf2fbe3467cc267 (router.bittorrent.com:6881): 648.218621ms
    ebff36697351ff4aec29cdbaabf2fbe3467cc267 (router.utorrent.com:6881): 873.864706ms
    2/2 responses (100.000000%)
