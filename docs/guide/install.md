---
title: 5 · Installing
parent: Part II · Using URLInsane
nav_order: 1
---

# Installing
{: .no_toc }

- TOC
{:toc}

URLInsane is a single static binary with no runtime dependencies. The reference
data — vocabulary, misspelling lists, keyboard layouts, the geolocation database
— is compiled into it, so there is nothing to fetch before the first scan.

## A release binary

Releases are published for a wide set of platforms: macOS (amd64, arm64), Linux
(amd64, arm64, 386, arm, riscv64), Windows (amd64, arm64, 386), FreeBSD, OpenBSD
and NetBSD. Each is accompanied by a `.sha512`.

```bash
# Linux amd64, adjust the version and platform to taste
VER=0.9.0
curl -LO https://github.com/rangertaha/urlinsane/releases/download/$VER/urlinsane-$VER-linux-amd64
curl -LO https://github.com/rangertaha/urlinsane/releases/download/$VER/urlinsane-$VER-linux-amd64.sha512
sha512sum -c urlinsane-$VER-linux-amd64.sha512
chmod +x urlinsane-$VER-linux-amd64
sudo mv urlinsane-$VER-linux-amd64 /usr/local/bin/urlinsane
```

Check the [releases page](https://github.com/rangertaha/urlinsane/releases) for
what is actually published — the platform matrix in the Makefile is what *can*
be built, and a given release may carry fewer.

## With Go

```bash
go install github.com/rangertaha/urlinsane/cmd/urlinsane@latest
```

Requires the Go toolchain, and puts the binary in `$(go env GOPATH)/bin`. This
is the quickest route if you already have Go, and the only route to an unreleased
commit.

## From source

```bash
git clone https://github.com/rangertaha/urlinsane.git
cd urlinsane
make build            # builds ./build/urlinsane and ./build/datasets
make install          # chmod +x and move it to /usr/local/bin (uses sudo)
```

Useful targets:

| Target | What it does |
|---|---|
| `make build` | build `urlinsane` and the `datasets` maintainer tool |
| `make test` | unit tests |
| `make race` | tests under the race detector — the run that catches a wrong graph |
| `make check` | `fmt` + `vet` + `race`, what CI runs |
| `make dataset` | rebuild the embedded `dataset.db` from `datasets/` |
| `make release` | cross-compile every platform, with checksums |
| `make dpkg` | build a Debian package |

Builds use `-trimpath` and strip symbols, so the same source produces a
byte-identical binary on another machine. Cross-compilation works for every
target above without a C toolchain because the SQLite driver is pure Go.

{: .note }
> `make dataset` writes `internal/config/dataset.db`, which `go:embed` picks up.
> It must run *before* `make build` for a data change to reach the binary.

## Verifying it works

```console
$ urlinsane --version
urlinsane version 0.9.0
```

```console
$ urlinsane typo --list formats
table
json
ndjson
csv
dot
```

That second command exercises the dataset path as well as the binary: on a first
run it creates `~/.config/urlinsane` and extracts the shipped data into it,
reporting what it did on stderr.

## What gets created on first run

```
~/.config/urlinsane/
├── dataset.db        vocabulary, misspellings, homoglyphs, source lists
├── maxmind.db.gz     geolocation database
├── blocks/           the content-addressed graph store (created on first --save-graph)
└── config.yaml       plugin settings (optional; nothing writes it for you)
```

`dataset.db` and `maxmind.db.gz` are embedded in the binary and written out only
**if absent**. That is worth knowing when you upgrade: an existing
`~/.config/urlinsane/dataset.db` is left alone, so a new binary with richer
reference data will keep using the old file. If `--list languages` shows fewer
languages than the release notes claim, that is why:

```bash
rm ~/.config/urlinsane/dataset.db     # it will be re-extracted on the next run
urlinsane typo --list languages
```

Failing to extract a file is not fatal. The operators that needed it are left
out of the compiled plan instead, and the run says so — a scan with no
geolocation must not silently look like a target with no geolocation.

## Uninstalling

```bash
sudo rm /usr/local/bin/urlinsane
rm -rf ~/.config/urlinsane
```

The second command deletes any saved scans in `blocks/` along with the reference
data.

---

Next: **[Your first scan](../first-scan/)**.
