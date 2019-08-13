# GoGet

## Description

GoGet installs Go (Golang) packages by executing `go get`.

## Settings

The following settings are available for the GoGet handler:

| yaml Tag | Type | Description |
|----------|------|-------------|
| `workingDir` | String | sets the directory in which the go get command will be executed. must exist |
| `updateDependencies` | Boolean | If true, execute the go get command with `-u`, updating the dependencies |
| `printCommandOutput` | Boolean | If true, print the go get command's output on the fly |

## Package Definition

The package name follows the Go convention with [URL]@[Version]. If no version
is set, it will use latest.

If packages are pinned to a specific version, the packages will not be upgraded.
If you want to upgrade pinned packages, use the install command with the new
version (or just set the version to "latest" anyway if you're as lazy as I am)
