---
title: Home
layout: home
nav_order: 1
description: >-
  URLInsane generates the names an attacker would have registered, then goes and
  looks to see whether they did.
permalink: /
---

# URLInsane

A book about typosquatting, and the tool that hunts it.
{: .fs-6 .fw-300 }

[Start reading](attack/){: .btn .btn-primary .mr-2 }
[Install](guide/install/){: .btn .mr-2 }
[GitHub](https://github.com/rangertaha/urlinsane){: .btn }

---

Someone types `exmple.com` instead of `example.com`. Someone installs a package
from memory and gets the one whose name is not quite it. Someone approves a pull
request from a GitHub account whose handle differs from a colleague's by one
transposed letter.

None of these are exploits. Nothing is bypassed, nothing overflows. The attack
is that a name a human meant to type and a name an attacker registered are close
enough that the human cannot tell them apart in the half second they spend
looking.

URLInsane enumerates that space. Give it a name and it generates the plausible
neighbours of that name — by keyboard geometry, by spelling, by sound, by
script, by the conventions of whichever registry the name lives in — and then
goes and observes which of them somebody already owns.

```console
$ urlinsane typo -a co -d 1 example.com
TYPE    KEY          DEPTH  EXISTENCE  RISK      DETAIL
domain  eample.com   0      live       high      created=2006-01-07 registrar=GoDaddy.com, LLC
domain  examle.com   0      live       high      created=2005-05-28 registrar=GoDaddy.com, LLC
domain  exampe.com   0      live       critical  created=2005-07-26 registrar=GoDaddy Online Services…
domain  exampl.com   0      live       high      created=2005-10-05 registrar=PDR Ltd.
domain  example.com  0      live                 created=1995-08-14 registrar=RESERVED-IANA
domain  exaple.com   0      live       critical  created=2005-07-27 registrar=GoDaddy.com, LLC
domain  exmple.com   0      live       critical  created=2004-12-08 registrar=PDR Ltd.
```

Every single-character deletion of `example.com` is registered, most of them two
decades ago, none of them by IANA. That is the whole subject in one screen.

## How this book is arranged

**[Part I — The attack](attack/)** is about the problem and has nothing to do
with the tool. What typosquatting is, who does it and why it pays, which named
things it happens to, and the reasons a human produces one string when they
meant another. Read this if you want to know what you are looking for.

**[Part II — Using URLInsane](guide/)** is the working manual. Installing,
scanning, choosing algorithms, controlling how far a scan reaches, reading the
report, and wiring the whole thing into CI so a build fails when a new lookalike
appears.

**[Part III — Inside the engine](internals/)** is for people changing the code.
URLInsane is a graph engine that happens to run a typosquatting workload; this
part explains what that means, why the pipeline it replaced had to go, and how
to add an operator without rewiring anything.

**[Reference](reference/)** is the flat material — the full CLI, the design
document, the keyboard model, a glossary, and the research the algorithms are
drawn from.

## What it works on

Typosquatting is not a domain problem that other ecosystems happen to share. It
is a *naming* problem, and every namespace humans type into has it. The kind of
target is detected from the string alone — there is no `--type` flag:

| You type | It reads |
|---|---|
| `urlinsane typo acme.com` | a domain |
| `urlinsane typo bob@acme.com` | an email — varies `bob`, `acme.com`, and the address |
| `urlinsane typo npm:lodash` | a package on a named registry |
| `urlinsane typo github.com/acme/tool` | a repository |
| `urlinsane typo bobsmith` | a username |

## Status

The engine is mid-rewrite, from a linear plugin pipeline to a graph engine, and
this book documents what the code does today rather than what the design
promises. Where something is specified but not yet wired up, it is marked:

{: .todo }
> Passages marked like this describe behaviour that is designed and not built.
> [CLI §9](CLI/#9-not-yet-implemented) tracks the full list.

## Credits

Written by [Rangertaha](https://github.com/rangertaha). Inspired by
[URLCrazy](https://morningstarsecurity.com/research/urlcrazy) and
[dnstwist](https://github.com/elceef/dnstwist), and by two decades of academic
measurement work credited in the [bibliography](reference/bibliography/).

Licensed under the GPLv3.
