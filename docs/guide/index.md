---
title: Part II · Using URLInsane
nav_order: 3
has_children: true
permalink: /guide/
---

# Part II · Using URLInsane

The working manual. Every command in this part was run against the binary that
built this book; where output is trimmed for width it says so, and where
behaviour is designed but not yet wired up it is marked with a red callout.

The shape of a scan never changes, whatever the target:

```
   target string
        │
        ├─ decompose ──▶ the entities the target is made of
        │                (an email is a local part, a domain, and an address)
        │
        ├─ vary ───────▶ plausible neighbours of each of those names
        │                (27 algorithms, keyboard- and language-aware)
        │
        ├─ observe ────▶ what actually exists out there
        │                (DNS, WHOIS, registries, reverse lookups)
        │
        └─ analyse ────▶ which findings are worth your attention
                         (risk, campaign clustering, dependency confusion)
```

Chapters:

1. **[Installing](install/)** — binaries, `go install`, building from source.
2. **[Your first scan](first-scan/)** — reading the default table, and the three
   flags that matter on day one.
3. **[Targets and scope](targets/)** — how a target string is classified, and
   how to narrow what gets varied without changing how the target is read.
4. **[Algorithms](algorithms/)** — all 27, what each one models, and how to pick
   a useful subset.
5. **[Observation and depth](observing/)** — operators, what `--depth` actually
   counts, and why "absent" and "unknown" are different answers.
6. **[Reading the report](reports/)** — filters, formats, risk levels, `--save`.
7. **[Automation](automation/)** — exit codes, `--fail-on`, JSON and Graphviz
   pipelines, running it in CI.
8. **[Datasets and languages](datasets/)** — where the vocabulary and keyboard
   data come from, and how to add a language.
