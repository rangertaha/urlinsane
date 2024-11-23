# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
- Removed the server
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

