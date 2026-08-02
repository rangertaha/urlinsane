---
title: 2 · The named-entity surface
parent: Part I · The attack
nav_order: 2
---

# The named-entity surface
{: .no_toc }

- TOC
{:toc}

## Domains are just the oldest case

Typosquatting is usually described as a domain problem. That is an accident of
history: domains were the first namespace that was globally shared, cheap to
enter, and typed by hand. Every namespace with those three properties has the
same problem, and most of the namespaces invented since are worse, because the
thing doing the typing is increasingly not a person reading carefully but a
machine following a manifest.

The general form: **anywhere a human writes a name that a machine resolves to a
thing, someone can register a name that is nearly it.**

URLInsane models five kinds of nameable entity. It works out which one you gave
it from the string alone — there is no `--type` flag, because a flag that
contradicted the string would just be a second way to be wrong:

| Type | Example | What makes it distinct |
|---|---|---|
| `domain` | `acme.com` | dotted name whose suffix is in the public suffix list |
| `email` | `bob@acme.com` | an `@` with a host after it |
| `package` | `npm:lodash` | `registry:name` |
| `repo` | `github.com/acme/tool` | `host/owner/name`, or any URL or scp-style remote of one |
| `username` | `bobsmith` | anything else that is a legal handle |

The order matters and is most-specific-first. Rule 4 is the interesting one:
`lodash` and `lodash.com` are both legal hostnames, so a hostname test alone
would classify every bare package name as a domain and go do DNS lookups for it.
Requiring a real public-suffix match is the only discriminator that does not
need a list of known package names.

Note also what is *refused*. Give it `192.0.2.1` and it declines the scan rather
than producing something: an address is a real entity but not a nameable one,
and varying it character by character produces addresses with no relationship to
the original. There is nothing there to typosquat.

## What decomposition buys

A target is rarely one name. `bob@acme.com` is three: the address itself, the
local part `bob`, and the domain `acme.com`. Each is separately squattable, and
the three attacks are different.

- Squat the **domain** and you receive misdirected mail for the whole company.
- Squat the **local part** — register `bobb@acme.com` if the platform allows,
  or the same handle on a different service — and you impersonate one person.
- Squat the **address** as a single string, `bob@acme-corp.com`, and you get the
  display-name attack that most mail clients still render as just "Bob".

URLInsane decomposes the target into its constituent entities and varies each,
which is why an email scan produces qualitatively different findings from a
domain scan of the same organisation. [Targets and scope](../../guide/targets/)
covers how to narrow that when you only want one of the three.

## Packages: the high-consequence case

A domain squat waits for a human to make a mistake. A package squat waits for a
*build* to make one, and builds run unattended, on a schedule, with credentials.

The mechanics differ from domains in ways that matter:

**Install-time code execution.** npm `postinstall` hooks, Python `setup.py`,
RubyGems extensions — several ecosystems run attacker-controlled code merely to
install a package, before anything is imported. The window between mistake and
compromise is one command.

**The victim is a machine.** A developer types a name once; CI installs it a
thousand times. And the name is often not typed at all but copied from a blog
post, a Stack Overflow answer, or an LLM's suggestion — none of which the reader
verifies against the registry.

**Namespaces are shallow and inconsistent.** npm has scopes (`@acme/tool`), Go
has full module paths, PyPI has a flat namespace where `-` and `_` are
interchangeable in some contexts and not others. Each inconsistency is a squat
opportunity: `python-dateutil` and `python_dateutil`, `acme-tool` and
`acme_tool`, `@acme/tool` and `acme-tool`.

**Dependency confusion.** The one that is not a typo at all. If a build resolves
a package name against both a private registry and a public one, and the public
one has no such package, then an attacker who registers that name publicly may
win the resolution — often because the public version number is higher. The
"typo" here is an organisation publishing its internal package names in a
manifest that leaks. URLInsane has an analyzer for this; see
[Analysis](../../internals/analysis/).

The tool ships existence checks for **13 package registries** — cocoapods,
conda, cpan, crates, hackage, hex, homebrew, npm, nuget, packagist, pub, pypi
and rubygems — plus **9 repository hosts** and **66 username platforms**. Those
lists are data in the dataset database, not code, so extending them is an import
rather than a release ([Datasets](../../guide/datasets/)).

## Repositories and usernames

A repository is a name in two parts, and both can be squatted independently:
`github.com/acme/tool` can be attacked as a different `acme` (a lookalike
account) or a different `tool` under the real account's lookalike.

The account-level attack is worth taking seriously because of what accounts
confer. A GitHub handle one transposition away from a maintainer's is a
plausible sender for a pull request, a plausible author on a commit, and in
ecosystems that resolve packages by repository path — Go, most obviously — the
handle *is* part of the package name. `github.com/acme/tool` and
`github.com/acmé/tool` resolve to different code, and the difference survives
copy-paste.

Usernames also cross platforms in a way domains do not. The same handle on
GitHub, npm, Docker Hub, and a company Slack reads to a human as the same
person, and there is no authority anywhere that makes that true. Registering a
handle on the one platform your target has not claimed yet costs nothing and is
indistinguishable from them having claimed it.

The algorithms that only apply to these types reflect their naming conventions
rather than keyboard geometry:

| Algorithm | Applies to | What it models |
|---|---|---|
| `afx` | package, repo, username | plausible prefixes and suffixes: `-js`, `node-`, `python-`, `-cli` |
| `nsc` | package, repo | moving a name between namespaces or scopes |
| `sep` | package, repo, username | swapping the separators a registry allows: `-` for `_` for `.` |

## Why one engine, not five tools

It would be possible to write a domain typosquatting tool, then a package one,
then a username one. The reason not to is that the interesting findings are the
ones that *cross* types.

A phishing campaign registers a domain, points MX at it, creates a matching
GitHub handle, and publishes a package under it. Each of those is a weak signal
alone. Together they are a campaign, and you can only see it if they live in the
same result set with edges between them.

That is the argument for the graph engine in [Part III](../../internals/): not
that a graph is a tidy way to store rows, but that a scan starting at one name
should be able to reach an entity of a different type, three hops away, and
still be one scan.

---

Next: **[Why names get mistyped](../errors/)** — the error taxonomy behind the
algorithm list.
