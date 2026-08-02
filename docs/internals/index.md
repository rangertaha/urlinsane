---
title: Part III · Inside the engine
nav_order: 4
has_children: true
permalink: /internals/
---

# Part III · Inside the engine

URLInsane used to be a pipeline: generate a list of domains, push the list
through a fixed sequence of collectors, print it. That design is gone. What
replaced it is a graph engine on which typosquatting is merely the first
workload.

The move is worth stating precisely, because "we added a DAG" is the wrong
summary and was the wrong design:

> The graph is the data, not the execution plan.

An earlier revision did have a DAG — of *plugins*, with declared dependencies
between them, Terraform style. That made plugin order load-bearing, which made
the cache unsound, which meant every new collector was a rewiring exercise. What
shipped instead is **binding by data**: an operator declares the pattern of
graph it consumes and the relations it emits, and the scheduler matches. Nobody
declares an order. Nothing is topologically sorted to decide what runs, because
the entity graph is cyclic by nature — `domain → ip → PTR → domain` is correct
data, not a bug.

Chapters:

1. **[The graph](graph/)** — nodes, edges, props, and why props are an ordered
   field list rather than a map.
2. **[Types and relations](types/)** — the registry, capabilities, edge classes,
   and what a "nameable" thing is.
3. **[Operators](operators/)** — binding, emission, the applier, and the merge
   rules that keep concurrent writes deterministic.
4. **[The plan and the scheduler](plan/)** — what `--explain` prints, and how
   rounds and barriers work at run time.
5. **[Limits and termination](limits/)** — depth accounting, the terminal
   variant rule, truncation, and failure as data.
6. **[Analysis](analysis/)** — belief, risk scoring, campaign clustering,
   dependency confusion.
7. **[Content addressing](addressing/)** — why two identical scans produce
   identical CIDs, and what that buys for diffing.
8. **[Extending URLInsane](extending/)** — adding an algorithm, an operator, a
   language, or a keyboard.

The full design document, longer and more formal than these chapters, is
[DESIGN]({{ site.baseurl }}/DESIGN/). Where the two disagree, the code wins and
the disagreement is a bug in whichever document is behind.
