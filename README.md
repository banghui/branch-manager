
# Branch Manager

![best branch manager](static/branch-manager.png)

CLI tool (For Macs) to help you manage your git branches easily

# Prerequisite
- git cli
- go

# Installation

1. Build binary

```
make build
```


2. Add the binary to `/usr/local/bin`

```
make link
```

Alternatively, you can choose to put the built binary, `bin/bm`, into any other directory that is included in your computer's `$PATH`

# Usage

```
> bm
```

# Uninstallation

Remove `bm` binary file from `/usr/local/bin`
