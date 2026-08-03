---
title: 10 · Reading the report
parent: Part II · Using URLInsane
nav_order: 6
---

# Reading the report
{: .no_toc }

- TOC
{:toc}

## Anatomy of the table

```
example.com                                   ← the target

TYPE    KEY          DEPTH  EXISTENCE  RISK    DETAIL
domain  eaxmple.com  0      live       high    created=2017-06-20 … registrar=GMO Internet Group…
domain  examlpe.com  0      live       medium  created=2016-07-11 … registrar=TurnCommerce, Inc.…
domain  example.com  0      live               created=1995-08-14 … registrar=RESERVED-IANA
tld     com          0      untried

FINDINGS                                      ← why the risk column says what it says
  HIGH    live-variant  eaxmple.com is live: [resolves accepts mail]
  MEDIUM  live-variant  examlpe.com is live: [resolves]

domain 37  ip 93  registrant 5  tld 2          ← node counts by type
137 nodes  184 edges  38 live  12 absent  0 unknown
4 rounds in 42.87s
```

| Column | |
|---|---|
| `TYPE` | node type — `domain`, `ip`, `tld`, `registrant`, `package`, `username`, … |
| `KEY` | the canonical name |
| `DEPTH` | observation hops from the seed. Variants are 0 |
| `EXISTENCE` | `live`, `absent`, `unknown`, `untried` — [three-state, deliberately](../observing/#three-state-existence) |
| `RISK` | assigned by analyzers over the finished graph |
| `DETAIL` | materialized props: WHOIS dates, registrar, TXT records, CNAME |

A `⏎` inside `DETAIL` is a line break in the underlying value — TXT records
frequently hold several strings, and they are joined for display rather than
truncated.

`--verbose` adds provenance and engine belief to every row:

```
domain  eaxmple.com  0  live  high  created=… decompose.domain:ok dns-a:ok dns-cname:empty
                                    dns-mx:ok dns-ns:ok dns-txt:ok idn:empty whois:ok belief=1.000
```

Each `operator:status` pair is one attempt. `ok` produced data, `empty` ran and
found nothing. This is how you tell "no MX record" from "the MX lookup failed" —
the distinction the existence column summarises.

## Filters select rows, never work

```bash
urlinsane typo -f live acme.com
urlinsane typo -f live -f 'depth<=0' acme.com     # both must hold
urlinsane typo -f 'risk>medium' acme.com
urlinsane typo -f type=domain acme.com
```

```
FILTER     SELECTS
live       an observation operator returned ok
absent     none did, and at least one determined absence
unknown    every attempt failed, timed out or was skipped
untried    no operator ran on it
risk>SEV   nodes with findings above a severity
type=NAME  a node type
depth<=N   observation hops from the seed
```

Two things to internalise:

**Filters apply to the report, not to the scan.** They select which rows are
printed after everything has been observed. Narrowing the *scan* is what
`--depth`, `--algorithm` and the scope positional do. This is why re-filtering a
saved scan never costs another lookup.

**`risk>SEV` is exclusive.** `risk>medium` returns `high` and `critical` rows,
not medium ones. Severities are `info`, `low`, `medium`, `high`, `critical`.

Multiple `--filter` flags are ANDed. The most useful combination in practice:

```bash
urlinsane typo -f live -f 'depth<=0' acme.com     # live variants only, no infrastructure
```

## Formats

| Format | Shape | Good for |
|---|---|---|
| `table` | aligned, coloured | reading |
| `json` | one document at end of scan | `jq`, storing a whole run |
| `ndjson` | one object per line | streaming, `grep`, log pipelines |
| `csv` | flat rows | spreadsheets, ticketing imports |
| `dot` | Graphviz | seeing the graph as a graph |

### json

```console
$ urlinsane typo -o json acme.com | jq '.nodes[] | select(.existence=="live") | .key'
```

```json
{
  "target": "example.com",
  "partial": false,
  "rounds": 4,
  "nodes": [
    {
      "id": "27b84379ec94aa8a",
      "type": "domain",
      "key": "eaxmple.com",
      "depth": 0,
      "existence": "live",
      "in_closure": true,
      "props": [
        { "name": "txt", "value": "v=spf1 mx -all", "kind": "string" }
      ],
      "statuses": [
        { "operator": "dns-a", "status": "ok" },
        { "operator": "whois", "status": "ok" }
      ]
    }
  ]
}
```

`"partial": true` means the scan was interrupted. Check it before treating an
empty result as good news.

### ndjson

One `{"kind":"run",…}` header object, then one `{"kind":"node",…}` per node:

```
{"kind":"run","partial":false,"rounds":4,"target":"example.com"}
{"kind":"node","type":"domain","key":"eaxmple.com","depth":0,"existence":"live",…}
```

### csv

```
type,key,depth,existence,risk,in_closure,props,findings,declined
domain,eaxmple.com,0,live,high,true,txt=v=spf1 mx -all,live-variant,
```

### dot

The graph, not a flattened list — nodes are coloured by existence and edges are
labelled with their relation:

```console
$ urlinsane typo -o dot acme.com > acme.gv
$ dot -Tsvg acme.gv > acme.svg
```

```
digraph urlinsane {
  rankdir=LR;
  node [shape=box style=rounded fontname="sans-serif"];
  "domain:eaxmple.com" [label="eaxmple.com\ndomain" style="rounded,filled" fillcolor="#ffe0b2"];
  …
}
```

This is the format that makes infrastructure clustering obvious: six variants
pointing at one nameserver is a picture, not a table.

## Saving

```bash
urlinsane typo acme.com --save report.csv       # format inferred from the extension
urlinsane typo acme.com --save report.json
urlinsane typo acme.com --save graph.gv         # dot
urlinsane typo acme.com --save report.txt       # table, uncoloured
```

Recognised extensions: `.json`, `.ndjson`, `.csv`, `.dot`/`.gv`, `.txt`/`.text`.
Anything else is an error rather than a guess:

```console
$ urlinsane typo acme.com --save out.weird
report: cannot infer a format from "out.weird"
```

A saved file is never coloured, whatever the terminal.

## Re-rendering a saved scan

`--save-graph` persists the graph itself — not a rendering of it — to the
content-addressed store, and prints its root CID:

```bash
urlinsane typo --save-graph acme.com
```

Afterwards, `report` renders it without re-observing anything:

```console
$ urlinsane report --scans example.com
WHEN              TYPE    TARGET       ROOT
2026-08-02 19:47  domain  example.com  bafyreiebberwkwmjjn6iucih6s4dk2nirosz37fffroiny3owepzs6qp6i
2026-08-02 19:39  domain  example.com  bafyreif2w4u4k5ksiymx3yyftbiin65sdhwescqdxq5u7tgo5v4b6ae6eu

$ urlinsane report -f live -f 'depth<=0' -o csv example.com
$ urlinsane report --at bafyreif2w4u4k5ksiymx3yyftbiin65sdhwescqdxq5u7tgo5v4b6ae6eu example.com
```

You name the **target**, not a file: `acme.com` is what you scanned and what you
remember; a root CID is not. `--at` addresses a specific earlier scan when you
do want the exact one.

Rendering never re-observes. The stored blocks are replayed through the same
applier the scan used and CID-checked against what was stored, so what you see
is byte-identical to what was saved. That is what makes "what changed since last
week" a comparison of two CIDs rather than two scans —
[Content addressing](../../internals/addressing/).

`report` takes the same `--filter`, `--output`, `--save` and `--verbose` flags as
`typo`, and one target argument only.

---

Next: **[Automation](../automation/)**.
