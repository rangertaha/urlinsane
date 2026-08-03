---
title: 21 · Learning what to scan
parent: Part III · Inside the engine
nav_order: 9
---

# Learning what to scan
{: .no_toc }

- TOC
{:toc}

## The problem

Generation is free and observation is not. A default scan of a domain produces
tens of thousands of candidates and each one costs a DNS query, a WHOIS lookup,
an HTTP request against somebody else's rate limit. When `--budget` or
`--frontier` binds, something has to decide which candidates survive.

Today nothing does. The engine's execution model is `uniformModel`:

```go
func (uniformModel) Initial() (float64, State)                 { return 1, nil }
func (uniformModel) Step(State, string, View) (float64, State) { return 1, nil }
```

Belief is 1 for everything, so expansion is unranked breadth-first and a bound
budget truncates alphabetically. `aaa-example.com` survives and `exmaple.com`
does not, for no reason at all.

## The shape of the answer

The seam has been there from the start. `graph.BeliefModel` is a hidden Markov
model over the **expansion tree**: a node's belief is its parent's latent state
pushed through the relation that admitted it, conditioned on the node's own
props. `internal/model` implements that — forward filtering, Baum-Welch,
Dirichlet smoothing, dag-cbor artifacts with CIDs — and `internal/eval` collects
ground truth from Certificate Transparency.

`internal/train` is the join. It featurizes a graph, walks it into paths, fits a
model, and orients it.

```
scan ──▶ --save-graph ──▶ store
                            │
                     Rehydrate + Finalize
                            │
                    Features ──▶ Paths ──▶ Fit ──▶ AnchorFocus
                                                        │
                                              graph.SetBeliefModel
```

## Train and serve must featurize identically

The hazard worth naming first, because it fails silently.

Training reads a finished graph. Inference reads a `graph.View` mid-scan,
restricted to what an operator declared it reads. Two code paths, one meaning —
and if they ever disagree, the model is fitted on symbols that never occur at
run time, belief collapses to the prior, no test goes red, and the scan simply
stops ranking.

So there is one feature function over the narrowest interface both can satisfy,
and `Trigger()` hands back the read-set declaration alongside it. A caller that
reads a field without declaring it is served a `View` that reports it unset.

That design caught a real bug immediately. The first featurizer emitted `algo:co`
— the algorithm that produced a variant, and the single most informative symbol
available. It is not obtainable: `VARIANT_OF` runs *origin → variant*, so on the
variant it is an **incoming** edge, and `View.Edges` returns outgoing edges only.

{: .todo }
> An operator mid-scan cannot see which algorithm produced the node it is
> looking at. That gap is one accessor wide — an incoming-edge view, declared in
> a trigger like any other read — and closing it would be the single biggest
> improvement available to this model.

## Paths, not rows

Baum-Welch is defined over sequences. The graph is cyclic (`domain → ip → PTR →
domain`) and has no sequences in it, so the sequences come from the expansion
*tree*: `Graph.Parent` gives every admitted node exactly one parent, fixed at a
barrier and never revised.

One path per leaf. An interior node appears in every path passing through it,
which is correct — it was observed once per descendant lineage.

{: .warning }
> **A rehydrated scan has no expansion tree.** Parents are assigned at a barrier
> and the store does not persist them: side tables carry depth and closure
> because those perturb nothing, but the tree is derived state. A graph replayed
> from the store arrives with every edge and no parents, and `Paths` refuses it.
> Call `Finalize` after loading — it runs a scheduler with no operators, which
> performs barrier 0 and stops.

## Three-state outcomes survive into the alphabet

`live`, `absent` and `unknown` are three emission symbols, not two.

Folding unknown into absent would train the model that a rate limit is evidence
a name is free — the one inference this codebase is arranged to prevent — and it
would do it invisibly, because the corpus would still look balanced. `Describe`
reports the outcome counts so that is checkable before trusting a model.

The same rule governs evaluation: only nodes with a settled observation are
scored. A node whose lookups all failed is not evidence either way, and counting
it as absent would score the model on the network's behaviour rather than on the
name's.

## Label switching, and why Focus is chosen after fitting

Baum-Welch is unsupervised. Which of its latent states a human would call
"promising" is whatever order EM converged to; the names in `Config.States` are
labels written on two buckets EM filled.

This is not a subtlety. Fitted on a real scan of `example.com`, the model
separated cleanly — 57 distinct belief values against uniform's 1 — and put
**every IPv6 address at 1.0 with every live typosquat at 0.0**. Ranking a
frontier by that is worse than not ranking it: the scan would spend its budget
walking address space and prune the squats.

`AnchorFocus` picks the state that actually emits `outcome=live`. One
comparison, the same tables, oriented on evidence. On the real scan it selects
the state *named* `dead` — which is the point.

```
before   1.0000 ip:2606:4700:58::adf5:3b7b   …   0.0000 domain:examle.com
after    1.0000 domain:exaple.com            …   0.0000 ip:2600:1f18:683a…
```

One implementation note: the alphabet check is membership in `Symbols()`, not
`SymbolIndex >= 0`. An unknown symbol maps to the out-of-vocabulary slot rather
than to `-1`, so the index test never fires and the model would be anchored on
"everything I have never seen".

## Does it work?

Trained on a saved scan of `example.com`, evaluated on a saved scan of
`paypal.com` that training never saw:

| | AUC |
|---|---|
| uniform | 0.500 |
| trained | **0.802** |

41 judged nodes — 35 live, 6 absent, 55 skipped as unknown or untried. The top
of the ranking is `payal.com`, `papal.com`, `aypal.com`, `paypl.com`,
`pypal.com`, `paypa.com`: every one live, every one a registered typosquat.

**Read that with its caveats.** One held-out scan and 41 judged nodes is an
encouraging signal, not a validated model. The base rate is 0.854, so
precision@k says almost nothing here — uniform already scores 1.000 at k=10 —
and AUC is the only informative number. The trained model's p@25 is in fact
slightly *worse*, 0.960 against 1.000, because it ranks the absent `tld:com`
node above a few live domains. That is a real weakness: the featurizer has no
symbol separating a structural node from a candidate beyond `type:`, and the
training corpus contained two `tld` nodes.

AUC counts ties as a half, which puts the uniform model at exactly 0.5 rather
than at whatever its tie-break order gives. That matters — 0.5 is the baseline,
and a model below it is doing harm rather than nothing. `auc()` is tested
against the inverted case, which must score 0.

## What is not built

{: .todo }
> **`--model` does not exist.** Nothing on the command line sets
> `scan.Options.Belief`, so every scan still runs the uniform model whatever has
> been trained. The wiring is small — load a model, `train.BeliefFrom(h)`,
> assign — and `CLI.md` §9 already lists `--model` among the specified flags.

Also open, in rough order of value:

1. **More scans.** 41 judged nodes is the binding constraint on every number
   above. Everything else is guesswork until there are dozens of scans to hold
   out properly.
2. **The incoming-edge accessor**, which would make `algo:` available and is the
   most informative symbol this model cannot see.
3. **A symbol separating structural nodes from candidates** — the `tld:com`
   mis-ranking is exactly that gap.

## Related

- [`pkg/aitypo`](https://pkg.go.dev/github.com/rangertaha/urlinsane/pkg/aitypo)
  is the other half of the AI work and a separate problem: it packages the
  generators as learnable *tasks* with exact oracles, for models that produce
  variants. This chapter is about learning which candidates are worth
  observing.
- [Analysis](../analysis/) explains why belief is engine-internal and analyzers
  cannot see it.
- [Limits](../limits/) covers the budgets and frontier that make ranking matter.
