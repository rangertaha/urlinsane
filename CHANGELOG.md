# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

The engine was rewritten from a linear plugin pipeline to a graph engine. Some
entries below describe the pipeline as it was being generalized; where the
rewrite superseded that work, the newer entry says so.

### Added
- **Graph engine.** Operators declare what data pattern they bind to and what
  they emit; the scheduler decides what runs when. The pipeline's `DependsOn`
  list made plugin order load-bearing and its cache unsound — binding by data
  means a new operator needs no rewiring, and `ptr` runs on any address
  whichever operator produced it. This replaces the dependency-ordered
  collector engine that preceded it.
- **Three-state existence**: live / absent / unknown. "We asked, it is not
  there" and "we could not tell" are opposite conclusions, and collapsing them
  turned a broken network into a clean bill of health.
- Generalized typosquatting to any named entity — domain, email, package, repo,
  username. The kind is detected from the target string; the `--type` flag it
  replaced is gone, and an optional scope positional narrows what gets varied.
- IPLD content-addressed graph store (`internal/store`): a dag-cbor filesystem
  blockstore, replay that rebuilds a graph through the applier and CID-checks it
  against what was stored, and cross-scan diffing. The root carries no timestamp
  or run id, so two identical scans address identically.
- Supply-chain detection for package/username/repo squatting:
  - `pkg`, `usr` and `repo` operators check a name against each source's
    existence API; the lists live in `datasets/sources/`.
  - Registry squat algorithms: separator substitution (`sep`), namespace
    confusion (`nsc`) and affix squatting (`afx`).
  - A dependency-confusion analyzer, which is only expressible because absence
    is recorded as absence rather than as failure.
- `datasets build` — rebuilds the embedded reference database from `datasets/`
  (`internal/dataset/gen`). Building deletes first rather than migrating: the
  shipped database carried columns from three schema generations, with the
  current ones empty because AutoMigrate can add a column but cannot fill it.
- Languages are seeded from the keyboard catalogue (`pkg/kb`), so a language the
  tool can reason about is listed whether or not its corpus is curated yet.
- Cross-platform release builds for 14 GOOS/GOARCH pairs, `-trimpath` and
  reproducible flags.
- `make check` — gofmt, vet and the race detector; plus `make dataset`, `make race`.

### Changed
- **SQLite driver is now pure Go** (`glebarez/sqlite` over `modernc.org/sqlite`).
  The cgo driver made cross-compilation need a C toolchain per target, and
  `CGO_ENABLED=0` produced a binary that ran but could not open the database —
  the failure was swallowed into an empty in-memory fallback, so every
  language-driven algorithm silently generated nothing.
- Go 1.26.5. CI pinned Go 1.22 while `go.mod` required 1.25, so it could not
  have built the tree; it now reads the version from `go.mod`.
- IPLD is the source of truth for scan results, superseding the GORM/SQLite
  results database. `dataset.db` remains SQLite, for reference data only.
- `internal/plugins` holds everything that acts on the graph, one directory per
  plugin, grouped by kind rather than by target. Each family is a library plus
  its plugins, with composition in `<family>/all`.
- Languages and keyboards stopped being plugins. A language is a dataset
  directory and a keyboard is a layout in `pkg/kb`; neither needs Go code.
- `--registered`/`--unregistered` became `--filter live|absent|unknown`, which
  also expresses the third case the two booleans could not.
- `--format`, `--file` and `--dir` collapsed into `-o` plus `--save`.
- Per-file licence headers reduced to a two-line SPDX identifier (182 files).

### Fixed
- `.gitattributes` marks binaries `-text`. `internal/config/maxmind.db.gz` had
  been corrupted by a text round trip — every byte >= 0x80 replaced with the
  UTF-8 replacement character, including the `8b` of its gzip magic, destroying
  11.7 MB of 49.7 MB. It has been unusable since commit 911a35c and still needs
  re-downloading; `scripts/mmdb.sh` now fetches and verifies it.
- `urlinsane typo` never called `config.Init()`, so no scan opened the dataset
  database: `--list languages` was empty and geo and the source operators were
  omitted from every plan.
- `make deps` ran `go get ./...`, which mutates `go.mod` in module mode.
- `make test` skipped `./datasets/...`.
- Self-heal a corrupt or unreadable `dataset.db` instead of panicking on startup.
- `--verbose` and `--debug` log to stderr rather than corrupting stdout.
- Removed the mandatory per-collector delay and the double plugin
  initialization per variant.

## [0.9.0] - 2025-01-07
### SQLite database backend
- Load typosquatting and linguistic datasets from the database
- Store results in the database for better analysis
- Cleaned up the plugins and reduced repeating code
- Updated net library

## [0.8.2] - 2024-11-17
### Fixed summary output
- Fixed some issues with getting the total record count and setting the live flag

## [0.8.1] - 2024-11-17
### Fixed img plugin
- Fixed a problem with the img plugin that was causing an issue.

## [0.8.0] - 2024-10-29
### Major rewrite of the engine and plugin system
- Using interfaces instead of functions for the plugins
- Redesigned the plugins to make it easier to author new ones
- Added support for arbitrary names, email addresses, and usernames
- Added an optional progress bar
- Improved output plugin that makes it easy to create new output formats
- Better documentation around the algorithms
- Publish a Debian package as well as binaries
- Caching results to improve performance
- Added combo squatting plugin that uses the keywords extracted from the text on the target domain
- Added an information plugin that provides the topics of each online domain
- VSM(Vector Space Model) plugin for comparing the similarity of two domains

## [0.7.0] - 2024-07-11
### Cleanup & reorganize codebase
- Removed the API server
- Removed the Dockerfile
- Improved plugins

## [0.6.1] - 2019-10-03
### Improvements
- Cleaned up the help output
- Added interface for an experimental storage backend

## [0.6.0] - 2019-07-30
### Improvements
- Improved performance
- Added additional data to each result record

## [0.5.2] - 2019-05-07
### API server
- Changed typo value

## [0.5.1] - 2019-05-07
### API server
- Changed endpoint method

## [0.5.0] - 2019-05-04
### Added API server
- Added API server
- Added 'typo' and 'server' commands.
- Improved structs for JSON output.

## [0.4.0] - 2019-01-17
### Added
- Code cleanup
- Fixed bug with NS record
- Fixed bug with keyboard selection
- Added the Armenian language
- Added French Canadian language

## [0.3.0] - 2018-09-10
### Added
- Added GeoIP function
- Updated help string
- I added a live filter for domains with IP addresses.
- Added the Armenian language.
- Added the Persian language.
- Added the Hebrew language.

## [0.2.0] - 2018-08-30
### Added
- Added concurrency to the extra functions
- Added SSDeep page similarity function
- Updated documentation
- Updated version and builds for Windows, Linux, OSX

## [0.1.0] - 2018-08-26
### Added initial code

