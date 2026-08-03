---
title: 11 · Automation
parent: Part II · Using URLInsane
nav_order: 7
---

# Automation
{: .no_toc }

- TOC
{:toc}

A typosquatting scan is worth little as a thing someone runs occasionally and
reads. It is worth a lot as a scheduled control whose output is a diff and whose
failure mode is a red build.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | clean |
| `1` | execution error |
| `2` | a finding at or above `--fail-on` |

That third code is what makes the tool a gate. Without it, the only way to react
to results is to parse stdout.

A second Ctrl-C, which aborts without writing a report, exits `130` — the shell
convention for a signal, and deliberately outside the documented `0`/`1`/`2`
contract, so a script that switches on those three never mistakes an abort for a
result.

```console
$ urlinsane typo -a cs -d 1 --fail-on high example.com > /dev/null; echo $?
2
$ urlinsane typo -a cs -d 1 --fail-on critical example.com > /dev/null; echo $?
0
```

`--fail-on` takes `info`, `low`, `medium`, `high` or `critical`, and triggers at
**or above** the level given — unlike the `risk>SEV` filter, which is exclusive.

## Gating a build on dependency squats

The case with the shortest fuse. A developer adds a dependency; you want to know
before it is installed on a machine with credentials.

```yaml
# .github/workflows/squat-check.yml
name: typosquat check
on:
  pull_request:
    paths: ['package.json', 'package-lock.json', 'requirements.txt']
  schedule:
    - cron: '0 6 * * 1'

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install urlinsane
        run: go install github.com/rangertaha/urlinsane/cmd/urlinsane@latest
      - name: Check our package names
        run: |
          urlinsane typo --fail-on high \
            -a afx,nsc,sep,cb,co,cs \
            --save-graph npm:acme-utils
```

Two notes on being a good citizen here: keep the algorithm list narrow so you
are not making thousands of registry requests on every pull request, and put the
broad scan on the schedule rather than on the PR trigger.

## Diffing: reporting only what is new

The point of `--save-graph` is that a scan is stored as content-addressed
blocks, so two identical scans produce the same root CID and a comparison is a
comparison of hashes, not a re-scan.

```bash
# Monday
urlinsane typo --save-graph acme.com

# Next Monday
urlinsane typo --save-graph acme.com
urlinsane report --scans acme.com
```

```
WHEN              TYPE    TARGET       ROOT
2026-08-02 19:47  domain  example.com  bafyreiebberwkwmjjn6iucih6s4dk2nirosz37fffroiny3owepzs6qp6i
2026-08-02 19:39  domain  example.com  bafyreif2w4u4k5ksiymx3yyftbiin65sdhwescqdxq5u7tgo5v4b6ae6eu
```

Same root CID means nothing changed — no diffing required, and no report worth
sending. Different roots mean something moved, and you can render either scan
with `report --at <cid>`.

{: .todo }
> A built-in `--resume` and graph diff are designed but not wired up. Today,
> comparing two runs means rendering both and diffing the output:
>
> ```bash
> urlinsane report --at $OLD -f live -o csv acme.com | sort > old.csv
> urlinsane report --at $NEW -f live -o csv acme.com | sort > new.csv
> diff old.csv new.csv
> ```

## Shipping findings somewhere

**To a ticket queue**, via CSV:

```bash
urlinsane typo -f live -f 'risk>medium' --save findings.csv acme.com
```

**To a log pipeline**, via NDJSON — one object per line, so it streams:

```bash
urlinsane typo -o ndjson acme.com | while read -r line; do
  echo "$line" | jq -c 'select(.kind=="node" and .existence=="live")'
done
```

**To a chat webhook**, via `jq`:

```bash
urlinsane typo -o json acme.com \
  | jq -r '.nodes[] | select(.existence=="live" and .depth==0) | "• \(.key)"' \
  | head -20
```

**As a picture**, via Graphviz — the format that makes shared infrastructure
obvious:

```bash
urlinsane typo -o dot acme.com | dot -Tsvg > acme.svg
```

## Cron

```cron
# Weekly brand sweep, saved for diffing, alert only on new high findings
0 6 * * 1  urlinsane typo --save-graph --fail-on high --workers 4 acme.com \
             --save /var/log/urlinsane/acme-$(date +\%F).csv \
             || mail -s "new typosquat findings: acme.com" secops@acme.com
```

Set `--workers` on anything scheduled. The default parallelism is tuned for an
interactive run, and a cron job that hammers WHOIS at 6am every Monday from the
same IP will get that IP blocked.

## Scripting notes

**Always check `partial`.** An interrupted or timed-out scan reports what it
collected and marks itself:

```bash
urlinsane typo -o json acme.com | jq -e '.partial == false' > /dev/null \
  || echo "warning: scan was cut short, results are incomplete"
```

**Do not treat `unknown` as `absent`.** In JSON, `existence` is one of `live`,
`absent`, `unknown`, `untried`. A script that tests `!= "live"` has just decided
that a rate-limited registry means nobody registered the name.

**Colour is off automatically when stdout is not a terminal**, and `--no-color`
forces it off. Saved files are never coloured.

**Pin the version.** `go install …@latest` in CI means your gate's behaviour
changes when the tool does. Pin a tag.

{: .todo }
> `--quiet` is documented in some places but is not implemented. Redirect stdout
> if you only want the exit code.

---

Next: **[Datasets and languages](../datasets/)**.
