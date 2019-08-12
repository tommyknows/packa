# Brew

## Description

Install formulae through `homebrew`.

## Settings

The following settings are available for the Brew handler:
| yaml Tag | Type | Description |
| `taps` | []string | a list of taps which to use as a source for formulae. The handler will automatically cleanup taps that are not listed here |
| `printCommandOutput` | Boolean | If true, print the go get command's output on the fly |

## formula Definition

formulae are defined through `[tap]/<formulaname>@[version]`. For example:

```sh
# Install formula "vim"
vim

# Install formula "vim" at version 8.1.0
vim@8.1.0

# Install formula "vim" from tap "mycool/tap":
mycool/tap/vim

# Intsall formula "vim" from tap "mycool/tap" at version 8.1.0:
mycool/tap/vim@8.1.0
```

If the version has been specified, the formula will be pinned automatically.

If a pinned formula is upgraded (through the `upgrade` command) with a new version,
the new version will be pinned again.
If a pinned formula is upgraded with no new version specified, it will be unpinned
and upgraded.

When upgrading all formulae, pinned ones will not be upgraded.

## Glossary

- Formula in brew is a package
- Formulae are thus multiple packages
