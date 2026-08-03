---
title: 9 · Observation and depth
parent: Part II · Using URLInsane
nav_order: 5
---

# Observation and depth
{: .no_toc }

- TOC
{:toc}

Generating a name costs nothing. Finding out whether anybody owns it costs a
network round trip against a service that is usually rate-limited and sometimes
run by volunteers. This chapter is about controlling that half of the scan.

## Operators

An **operator** is anything that expands or observes the graph. Where the old
design ran collectors in a fixed order over a list of domains, an operator
declares what data pattern it *binds to* and what it *emits*, and the scheduler
decides what runs when.

```bash
urlinsane typo --list operators
```

```
ID                 VERSION  RESOURCE  BINDS ON               EMITS
decompose.domain   1                  domain                 TLD_OF
decompose.email    1                  email                  LOCAL_PART,DOMAIN_OF
decompose.package  1                  package                OWNER
decompose.repo     1                  repo                   HOSTED_ON,OWNER
dns-a              1        dns       domain                 RESOLVES_TO
dns-cname          1        dns       domain
dns-mx             1        dns       domain                 MX
dns-ns             1        dns       domain                 NS
dns-txt            1        dns       domain
idn                1                  domain
pkg                1        http      package                EXISTS_ON
ptr                1        dns       ip                     PTR_TO
repo               1        http      repo                   EXISTS_ON
usr                1        http      username               EXISTS_ON
whois              1        whois     domain                 REGISTERED_BY
```

plus one operator per algorithm, all emitting `VARIANT_OF`.

**Binding is by data, not by producer.** `ptr` binds to any `ip`, so it runs on
addresses whether `dns-a` found them or something else did. That is what lets a
new operator slot in without anyone rewiring an order, and it is the central
idea of [Part III](../../internals/operators/).

The `RESOURCE` column names the shared, rate-limited thing an operator consumes
— `dns`, `whois`, `http`. It is how the scheduler avoids pointing forty
concurrent workers at one WHOIS server.

Operators whose data is missing are **omitted from the plan** rather than
failing at run time. If `geo` never appears, its geolocation database did not
load; if `pkg`, `usr` and `repo` never appear, the source lists are not in the
dataset. `--explain` shows you exactly which operators a run will use.

## Depth

```bash
urlinsane typo -d 1 acme.com
```

`--depth` counts **observation hops from the seed**, and only observation hops.

| Edge class | Relations | Depth cost |
|---|---|---|
| Structural | `TLD_OF`, `LOCAL_PART`, `DOMAIN_OF`, `OWNER`, `HOSTED_ON` | 0 |
| Variant | `VARIANT_OF` | 0 |
| Observation | `RESOLVES_TO`, `NS`, `MX`, `PTR_TO`, `REGISTERED_BY`, `EXISTS_ON`, `IN_ASN`, `MANIFEST` | 1 |

The dividing line is **whether the edge required a network call**. Decomposition
is derived from parsing the target string, so it is free; a variant is generated
locally, so it is free too. Everything that had to ask the network costs one.

So at depth 1:

```
acme.com  (depth 0)
   ├── VARIANT_OF ──▶ acmee.com     depth 0   ← variants are free
   └── RESOLVES_TO ─▶ 192.0.2.7     depth 1   ← one hop
                          └── PTR_TO ──▶ host.example.net   depth 2 — beyond the limit
```

Which is exactly what you see in a real report: variants at depth 0, their
nameservers and addresses at depth 1, and whatever those resolve to marked
`untried` at depth 2.

{: .warning }
> **`-d 0` means unbounded, not zero hops.** The limit is only applied when it
> is greater than zero. To *see* only the variants, filter the report instead:
> `--filter depth<=0`.

Practical settings:

| | |
|---|---|
| `-d 1` | the variants and their immediate records. Fast, and enough for monitoring |
| `-d 2` | adds what those records point at — reverse lookups, shared hosting |
| `-d 3` | the default; enough to see infrastructure clusters |
| `-d 0` | no limit. Use with a budget, or expect a long walk through address space |

Depth is what stops the walk, but it is not the only thing: see
[Limits](../../internals/limits/) for the round cap, the frontier cap and the
terminal-variant rule that keeps variants from being varied recursively.

## Three-state existence

The most important design decision in the observation half of the tool:

| State | Means |
|---|---|
| `live` | an observation operator returned ok |
| `absent` | none did, and at least one **determined** absence |
| `unknown` | every attempt failed, timed out, or was skipped |
| `untried` | no operator ran on it — usually past `--depth` |

"We asked, it is not there" and "we could not tell" are opposite conclusions.
Collapsing them — which most tools do, by treating any non-answer as *not
registered* — turns a rate limit, a broken resolver or a firewall into a clean
bill of health. That is the failure mode that gets someone breached while
holding a green report.

The summary line counts all of them:

```
137 nodes  184 edges  38 live  12 absent  0 unknown
```

A scan with a large `unknown` count is not a scan with good news. It is a scan
you should run again.

## Bounding the cost

The flags below are registered but hidden from the common help, because most
runs never touch them. They are the ones to reach for when a scan is too slow,
too large, or too rude to a registry.

| Flag | Controls |
|---|---|
| `--rounds` | backstop for a type flow that never converges (default 64) |
| `--workers` | concurrent operator calls |
| `--budget` | global admitted-node cap; `0` is unbounded |
| `--frontier` | cap on nodes admitted per round; declined candidates go to the ledger |
| `--attempts` | per-pair attempts within a round |
| `--timeout` | bound on a single operator call |
| `--no-color` | disable ANSI styling |

```bash
urlinsane typo --workers 4 --timeout 5s --budget 2000 acme.com
```

{: .warning }
> An unrestricted scan makes thousands of DNS, WHOIS and HTTP requests to
> infrastructure you do not own. WHOIS servers rate-limit aggressively and will
> block you; package registries publish acceptable-use policies. The defaults
> exist because the unbounded version of this tool is a load generator pointed
> at services you depend on.

## Interrupting

Ctrl-C stops expansion at the **end of the current round**. The barrier still
runs, so parents, belief and the truncation ledger are finalised rather than
left half-computed; analyzers then run over what exists and the report is marked
`PARTIAL` — in every format, including JSON, so a truncated scan is never
mistaken downstream for a complete one. A second Ctrl-C aborts without a report.

---

Next: **[Reading the report](../reports/)**.
