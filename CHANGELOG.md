# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Dependency-aware DAG collector engine: collectors run in dependency order
  (derived from each collector's declared dependencies), parallel across
  variants and sequential within a variant; `context` is propagated end to end
  with real per-collector timeouts.
- Generalized typosquatting to any named **entity** via `--type`
  (`domain`, `name`, `user`, `package`), including automatic classification of
  the target; all processing plugins (algorithms, collectors, analyzers)
  declare, and are filtered by, the entity types they support.
- IPLD content-addressed result store (`internal/store`): a dag-cbor filesystem
  blockstore with an IPLD schema (keyed union over entity types), a SQLite
  secondary index for name lookups, and cross-scan diffing.
- Supply-chain detection for package/username/repo squatting:
  - `pkg` (13 package registries), `usr` (66 username platforms) and `repo`
    (9 git forges) collectors that check each variant against a source's
    existence API; reference lists live in `datasets/sources/`.
  - `--manifest FILE` parses a project's declared dependencies
    (requirements.txt, package.json, go.mod, Cargo.toml, pyproject.toml,
    Gemfile, composer.json) and typosquat-checks each one.
  - Package-registry squat algorithms: separator substitution (`sep`),
    namespace/scope confusion (`nsc`) and affix combosquatting (`afx`).
  - `dc` dependency-confusion analyzer: flags names that exist on some
    registries but remain available (squattable) on others.

### Changed
- IPLD is now the source of truth for scan results, superseding the GORM/SQLite
  results database (the `dataset.db` reference data stays on SQLite).
- The Analyzers stage now runs with origin→variant pairing.
- Cleaner JSON output (dropped persistence/ID noise from records).

### Fixed
- Self-heal a corrupt or unreadable `dataset.db` instead of panicking on startup.
- Stop duplicating collected DNS/IP records across collectors and on re-scans.
- `--format <invalid>` no longer panics; it fails fast with a clear error.
- `--file`/`-o` now actually writes output (valid JSONL for the `json` format).
- `--verbose` now enables logging, and `--debug`/`--verbose` logs go to stderr
  instead of corrupting the results on stdout.
- Removed the mandatory 5-second-per-collector delay and the per-variant double
  plugin initialization.
- Corrected the `--distance` description (it is the *maximum* Levenshtein
  distance, not the minimum).

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

