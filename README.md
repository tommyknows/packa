# Packa

[![GoDoc](https://godoc.org/github.com/tommyknows/packa?status.svg)](https://godoc.org/github.com/tommyknows/packa)

A (meta) package manager with a declarative API.
Currently supporting the following package managers:

- go (go get)
- brew

For adding more managers, see [this](pkg/handlers/README.md)

## Prerequisites

A working `go` environment, preferably with the newest version of go (there have
been various bugs up until at least 1.12.5).

Make sure to set `GO111MODULE` to `on`.

## Installation

Initial installation happens with go get:

```
go get github.com/tommyknows/packa
```

After that, run

```
packa upgrade
```

to write and add packa itself to the initial config file.

## Usage

See `packa -h`.
