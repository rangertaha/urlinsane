---
title: 4 · What defenders do about it
parent: Part I · The attack
nav_order: 4
---

# What defenders do about it
{: .no_toc }

- TOC
{:toc}

A variant list is not a result. It is an input to a decision, and which decision
depends on who you are. This chapter is about what the list is *for*, and where
it stops being useful — which is worth knowing before you generate ten thousand
names.

## The four responses

**Register it.** Defensive registration works and does not scale. Buy the
handful of variants that are both cheap and high-yield: the single-character
deletions, the adjacent-key substitutions of the first three characters, the
`.net`/`.org`/`.co` of your `.com`, and the plural. Then stop, because the
neighbourhood is unbounded once other scripts and TLDs are admitted, and the
budget is not. Registration is prevention for the top 20 candidates and nothing
else.

**Monitor it.** The scalable response. Generate the neighbourhood, record what
exists today, and look at what changed. A squat that has been parked at a
registrar since 2005 is background noise; the same name acquiring MX records
last Tuesday is an event. The value is entirely in the delta, which is why
URLInsane addresses its results by content — two identical scans produce
identical CIDs and a diff is a comparison, not a re-scan
([Content addressing](../../internals/addressing/)).

**Block it.** Feed the variant list to the controls that can act on a name:
resolver blocklists, mail gateway rules, proxy categories, EDR. This is the
response with the best cost-to-benefit ratio and the least glamour. Note that
you can block variants you have *not* observed to exist — pre-emptive blocking
costs nothing per name and closes the window between registration and detection.

**Take it down.** UDRP, registrar abuse reports, platform reports for package
and account squats. Slow, per-item, and worth reserving for the cases with real
harm behind them: active phishing, malware, or a package in your dependency
tree. Registries and registrars respond to evidence of abuse far faster than to
assertions of trademark.

## Reading a result list

Most of a variant list is not interesting, and knowing what makes a row
interesting is most of the skill.

**Existence is the first cut, and it is three-valued.** *Live* means an
observation succeeded. *Absent* means at least one operator determined the name
is not there. *Unknown* means every attempt failed, timed out, or was skipped.
Collapsing the last two — as most tools do — turns a broken network or a
rate-limited registry into a clean bill of health, which is the single most
dangerous failure mode a tool like this has.

**Then: is it configured for anything?** A registered name that resolves
nowhere is speculation or parking. A name with MX records is being used, or is
waiting to receive mail that is not meant for it. A name with a certificate and
an A record pointing at a host serving your logo is an incident.

**Then: who else does it look like?** Squats cluster. The same registrar, the
same nameservers, the same hosting provider, the same registration date, the
same WHOIS privacy service. When six of your variants share infrastructure, you
are looking at one actor with a list, not six opportunists — and that changes
both the priority and the response, because the actor's *other* targets are
findable from the same infrastructure.

**Then: how old is it?** Registration date is the cheapest signal in the data.
Domains registered before your brand existed are usually unrelated. Domains
registered within a week of a product announcement are not.

## What a variant list cannot tell you

Three honest limits.

**It cannot tell you intent.** `exampl.com` registered by a domain investor in
2005 and `exampl.com` registered by a phishing crew last month look identical in
DNS. Attribution needs content, hosting history and behaviour — all outside what
name enumeration can see.

**It cannot be complete.** The candidate space is infinite once you admit
arbitrary combosquats, every TLD, and every script. Any tool that implies
completeness is lying about the algorithm set it ran. What the tool can honestly
say is which generators it ran and where it truncated — which is why URLInsane
keeps a truncation ledger rather than silently capping output
([Limits](../../internals/limits/)).

**It cannot distinguish "not registered" from "not visible to you".** Rate
limits, geo-fenced DNS, registries that require authentication, and providers
that answer differently to unfamiliar clients all produce absence that is not
absence. Hence the three-state existence model, again.

## Where this fits operationally

For a security team, a typosquatting scan is a **periodic control**, not an
investigation tool. The useful shape is:

1. A scheduled scan per brand-critical name, saved with `--save-graph`.
2. A diff against the previous scan, so the output is "what is new" rather than
   "everything".
3. New live findings above a severity go to the queue that already handles
   phishing reports.
4. The full variant list — live or not — goes to the blocking controls.

For a software team the shape is different, and tighter: run it against your
package names on a schedule, and against your *dependency* names on every
manifest change, with `--fail-on` set so the build stops rather than emailing
someone. Dependency confusion is the case where the gap between detection and
compromise is a single `npm install`, and a CI gate is the only control fast
enough. [Automation](../../guide/automation/) covers the wiring.

## A note on scanning responsibly

Generating names is free and harmless. Observing them is not entirely.

An unrestricted scan makes thousands of DNS queries, WHOIS lookups and HTTP
requests to registries that are, in most cases, run by volunteers or funded by
nobody in particular. WHOIS servers rate-limit aggressively and will block you.
Package registries publish acceptable-use policies. The defaults in URLInsane —
bounded depth, per-round frontier caps, a worker limit, per-call timeouts —
exist because the unbounded version of this tool is a denial-of-service tool
pointed at infrastructure you depend on.

Two rules that are not negotiable: do not use the output to register names that
infringe on others, and do not point the observation side at a service whose
terms you have not read. The tool makes it very easy to be a bad citizen at
scale.

---

That is Part I. **[Part II](../../guide/)** is the manual.
