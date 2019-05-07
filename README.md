# Packa

A package manager for go with a declarative API.

## Wait, Package Manager?

Yes. Not as go modules for dependencies, but more as `apt`
for go packages.

## Prerequisites

Make sure to set `GO111MODULE` to `on`.

Also, your go version needs to be _at least_ 1.12!

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
