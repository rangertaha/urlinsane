---
title: 6 · Your first scan
parent: Part II · Using URLInsane
nav_order: 2
---

# Your first scan
{: .no_toc }

- TOC
{:toc}

## The whole thing in one command

```bash
urlinsane typo example.com
```

That generates variants with every registered algorithm and observes each one to
a depth of three. It is the right command for finding out what the tool does and
the wrong command for a first run against a real target, because it is slow and
noisy. Start narrower:

```bash
urlinsane typo -a cs -d 1 example.com
```

`-a cs` runs one algorithm — character swapping, i.e. transpositions — and
`-d 1` stops observation one hop from the seed. Here is the actual output,
with long WHOIS detail trimmed to fit:

```console
$ urlinsane typo -a cs -d 1 example.com
example.com

TYPE    KEY                     DEPTH  EXISTENCE  RISK    DETAIL
domain  amir.mx.cloudflare.net  1      live
domain  ara.ns.cloudflare.com   1      live
domain  delta.netnautics.net    1      live               registrar=Tucows Domains Inc.
domain  dns1.eaxmple.com        1      live
domain  dns2.eaxmple.com        1      live
domain  eaxmple.com             0      live       high    created=2017-06-20 registrar=GMO Internet Group…
domain  examlpe.com             0      live       medium  created=2016-07-11 registrar=TurnCommerce, Inc.…
domain  exampel.com             0      live       high    created=2005-04-13 registrar=DNC Holdings, Inc.
domain  example.com             0      live               created=1995-08-14 registrar=RESERVED-IANA
domain  exapmle.com             0      live       high    created=2014-12-23 registrar=Dynadot Inc
domain  exmaple.com             0      live       high    created=2004-10-08 registrar=Cloudflare, Inc.
…

FINDINGS
  HIGH    live-variant  eaxmple.com is live: [resolves accepts mail]
  HIGH    live-variant  exampel.com is live: [resolves accepts mail]
  HIGH    live-variant  exapmle.com is live: [resolves accepts mail]
  HIGH    live-variant  exmaple.com is live: [resolves accepts mail]
  MEDIUM  live-variant  examlpe.com is live: [resolves]
  MEDIUM  live-variant  xeample.com is live: [accepts mail]

domain 37  ip 93  registrant 5  tld 2
137 nodes  184 edges  38 live  12 absent  0 unknown
4 rounds in 42.87s
```

Everything worth knowing about the tool is visible in that screen.

## Reading it

**The rows are the graph, not a list of variants.** `exmaple.com` is a variant
of the seed and sits at depth 0. `dns1.eaxmple.com` is not a variant of anything
— it is a nameserver the scan found *because* it looked up `eaxmple.com`, and it
sits at depth 1. Both are nodes; the report prints all of them, and `DEPTH` is
how you tell them apart. So is `--filter depth<=0`, if you only want the
variants.

**`TYPE` is not always `domain`.** The summary line shows what this scan
produced: 37 domains, 93 IP addresses, 5 registrants and 2 TLDs. A scan of an
email or a package produces other types. Filter with `--filter type=domain`.

**`EXISTENCE` has three values, and the summary counts all three.** `38 live 12
absent 0 unknown` — *live* means an observation succeeded, *absent* means an
operator determined the name is not there, and *unknown* means every attempt
failed, timed out, or was skipped. A row can also be *untried*, which means no
operator ran on it at all, usually because it sat past `--depth`. The
distinction between *absent* and *unknown* is deliberate and load-bearing; see
[Observation and depth](../observing/).

**`RISK` is a judgement, `FINDINGS` is the explanation.** The risk column is
assigned by analyzers that run over the finished graph, and the findings block
tells you what drove it. `[resolves accepts mail]` on a transposition of your
domain is a different problem from `[resolves]` — the first can receive
misdirected email.

**The seed itself appears, unflagged.** `example.com` is in the table at depth 0
with no risk. That is deliberate: you want to see the real thing's records next
to its neighbours' so you can compare registrars, nameservers and dates.

**The last line is the cost.** Four rounds, 42.87 seconds, for *one* algorithm
at depth 1. All 27 algorithms at the default depth of 3 is a different order of
magnitude. Generation is microseconds; the time is entirely network.

## Three flags for day one

| Flag | Why |
|---|---|
| `-a, --algorithm` | Restrict generation. `-a cs,co,acs` is a reasonable quick pass; `-a ^bf` runs everything except bit flipping |
| `-d, --depth` | How far observation walks from the seed. Default 3. **`-d 0` means no limit**, not "no hops" — use `--filter depth<=0` to see only the variants |
| `-o, --output` | `table` (default), `json`, `ndjson`, `csv`, `dot` |

And one for finding out what a build can do:

```bash
urlinsane typo --list algorithms   # 27 of them
urlinsane typo --list operators    # what expands and observes the graph
urlinsane typo --list languages    # what the dataset carries
urlinsane typo --list keyboards    # layouts available for adjacency
urlinsane typo --list types        # node types and their capabilities
urlinsane typo --list relations    # edge types, class and depth cost
urlinsane typo --list filters      # what --filter accepts
urlinsane typo --list formats      # what --output accepts
```

## Before you run: what will this do?

`--explain` compiles the plan and prints it without running anything. It costs
nothing and it is the honest answer to "what is this command about to do to the
network":

```console
$ urlinsane typo --explain example.com
plan 7437e47542eddfee
seed domain example.com
scope every nameable node in the seed closure
limits depth=3 rounds=64

type flow
  0  domain, email, ip, package, repo, username  (cycle)
  1  platform
  1  registrant
  1  tld

operators
  aci            on nameable           where in-seed-closure
  acs            on nameable           where in-seed-closure
  …
  decompose.domain on domain
  dns-a          on domain             [dns]
  dns-cname      on domain             [dns]
  dns-mx         on domain             [dns]
  dns-ns         on domain             [dns]
  dns-txt        on domain             [dns]
  idn            on domain
  ptr            on ip                 [dns]
  whois          on domain             [whois]
```

Three things to read here:

- **`plan <hash>`** identifies the compiled plan. The same flags against the
  same build produce the same hash.
- **The type flow** is the shape of what can be reached from this seed. The
  `(cycle)` marker is not a warning: `domain → ip → PTR → domain` is a genuine
  cycle in the data, and the engine does not rely on acyclicity to terminate.
- **`[dns]`, `[whois]`, `[http]`** are resource tags — which shared, rate-limited
  thing an operator consumes.

Operators whose data is missing are simply absent from the plan rather than
failing at run time. If you never see `geo`, its geolocation database did not
load; if you never see `pkg`/`usr`/`repo`, the source lists are not there.

## When something goes wrong

**`cannot tell what "…" is`** — target detection failed. See
[Targets and scope](../targets/) for the rules; the usual cause is a bare
package name with no registry prefix (`lodash` rather than `npm:lodash`).

**`… is an IP address; there is nothing to typosquat`** — deliberate. Scan the
name that resolves to it.

**`geolocation unavailable: …`** on stderr — the geo database did not open, the
`geo` operator was dropped from the plan, and the scan continued. Everything
else still ran.

**It hangs on a big scan.** It probably is not hanging; WHOIS and registry
lookups are slow and rate-limited. `--verbose` shows what is happening. Ctrl-C
is safe — an interrupt stops expansion and still reports what was collected,
with an exit code that says the scan was cut short.

---

Next: **[Targets and scope](../targets/)**.
