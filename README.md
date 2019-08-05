# Packa

[![GoDoc](https://godoc.org/git.ramonruettimann.ml/ramon/packa?status.svg)](https://godoc.org/git.ramonruettimann.ml/ramon/packa)

A (meta) package manager with a declarative API.
Currently supporting the following package managers:

- go (go get)

For adding more managers, see [this](pkg/handlers/README.md)

## Prerequisites

A working `go` environment, preferably with the newest version of go (there have
been various bugs up until at least 1.12.5).

Make sure to set `GO111MODULE` to `on`.

## Installation

Initial installation happens with go get:

```
go get git.ramonruettimann.ml/ramon/packa
```

After that, run

```
packa upgrade
```

to write and add packa itself to the initial config file.

## Usage

See `packa -h`.
