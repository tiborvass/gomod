# gomod

## Why?

This is a standalone temporary fork of the [go mod](https://golang.org/cmd/go/#hdr-Module_maintenance) command,
that has full support for [vndr](https://github.com/LK4D4/vndr)'s vendor.conf format.

As soon as https://github.com/golang/go/issues/25556 is resolved, I'll archive this repository.

## Install

```
$ go get -u github.com/tiborvass/gomod
```

## Usage

Identical to `go mod`.

## Build

```
$ pushd /path/to/github.com/tiborvass/go     # fork of github.com/golang/go
$ git checkout gomod && git pull
$ popd
$ GOROOT=/path/to/github.com/tiborvass/go ./generate.bash
$ go install
```
