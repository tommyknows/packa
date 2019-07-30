# Packa

A (meta) package manager with a declarative API.
Currently supporting the following package managers:

- go (go get)

## Prerequisites

Make sure to set `GO111MODULE` to `on`.

Also, your go version should be as new as possible to work with the goget
handler, as there have been various bugs up until at least 1.12.5.

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

## GoDoc

[![GoDoc](https://godoc.org/git.ramonruettimann.ml/ramon/packa?status.svg)](https://godoc.org/git.ramonruettimann.ml/ramon/packa)
