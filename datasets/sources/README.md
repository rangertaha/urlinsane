# Entity sources

Where to look up each kind of entity — the platforms, providers, and registries
that a generated variant can be checked against. These power lookups for the
non-domain entity types (`--type user|package`, emails) the same way the public
suffix list powers domains.

One source per line.

| File | Format | Used for |
|------|--------|----------|
| `usernames.lst` | `id url_template` | username/handle squatting across social & dev platforms |
| `packages.lst`  | `id url_template` | package/dependency squatting across registries |
| `email.lst`     | `provider_domain` | common email providers (one domain per line) |

In the URL templates, `%s` is replaced with the entity value to form a lookup
URL — e.g. `github https://github.com/%s` → `https://github.com/rangertaha`.

These lists are curated starting points (the most common squatting surfaces),
not exhaustive — extend them as needed.
