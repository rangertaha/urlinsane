# Package name corpus

Curated lists of well-known package names per code registry, used as a corpus /
seed target set for package (supply-chain / dependency-confusion) typosquatting.

One package name per line; one file per registry:

| File | Registry |
|------|----------|
| `pypi.lst`     | PyPI (Python) |
| `npm.lst`      | npm (JavaScript) |
| `crates.lst`   | crates.io (Rust) |
| `rubygems.lst` | RubyGems (Ruby) |

These are common typosquatting targets, not exhaustive snapshots of each
registry — extend as needed. The registries themselves (lookup URLs) live in
`datasets/sources/packages.lst`.
