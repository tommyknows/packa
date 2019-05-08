# Packa

A package manager for go with a declarative API.

## Wait, Package Manager?

Yes. Not as go modules for dependencies, but more as `apt`
for go packages.

## Prerequisites

Make sure to set `GO111MODULE` to `on`.

Also, your go version needs to be _at least_ 1.12! (see Issues for more info)

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

Duh, use `packa -h`.

## Issues

Current issues are:

### Versioning packa configs

Having the `.packa` directory checked in in a git repository results in the
error

```
📦 could not install package(s): Encountered error(s) while handling packages:
golang.org/x/tools/cmd/gopls@latest: could not install package golang.org/x/tools/cmd/gopls: go: cannot determine module path for source directory /Users/ramon/Documents/Tools/Dotfiles (outside GOPATH, no import comments)
exit status 1
```

This is because `go get` has an issue as documented [here](https://github.com/golang/go/issues/30515#issuecomment-490480624)
and is fixed at head of go. This means it will be fixed with the next release
of go (`1.12.6` or `1.13`)


