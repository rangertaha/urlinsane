---
title: 1 · What typosquatting is
parent: Part I · The attack
nav_order: 1
---

# What typosquatting is
{: .no_toc }

- TOC
{:toc}

## The mechanism

A typosquat is a name registered because it is *nearly* someone else's name.

That is the entire mechanism. There is no vulnerability. The victim's software
works perfectly: they typed `exmaple.com`, DNS resolved `exmaple.com`, the
browser connected to `exmaple.com`. Every layer did exactly what it was told.
The failure is one layer up, in the human, and it happened before the first
packet.

This is why typosquatting is durable in a way that technical vulnerabilities are
not. You cannot patch it. The set of names close to your name is not a property
of your infrastructure; it is a property of your name, your users' keyboards,
their languages, and the rules of whatever registry you are in. It changes when
you rebrand, and otherwise never.

What it costs an attacker is a registration fee. What it costs the defender is
the entire space of near-names, which for a seven-letter domain in a single TLD
already runs to thousands of candidates before you consider other TLDs, other
scripts, or other registries.

## The asymmetry

The defender's problem is combinatorial and the attacker's is not.

An attacker needs one name that one person mistypes. They do not need to guess
which mistake you will make, because they can wait: a squat registered today
collects traffic for years, and the cost of holding it is a renewal fee against
whatever the traffic monetises. Registration is cheap, mistakes are constant,
and the yield does not have to be high to clear the bar.

A defender who wants to *prevent* the attack has to enumerate and register the
whole neighbourhood, which nobody does past the obvious handful, because the
neighbourhood is unbounded once you admit other TLDs and other scripts. So
defence in practice is not prevention. It is **enumeration and observation** —
know what your neighbourhood looks like, know who is in it, and notice when that
changes. That is the job this tool exists to do.

## Why it pays

The uses fall into a small number of business models, and knowing which one you
are looking at changes what you should do about a finding.

**Parking and ad revenue.** The oldest and still the most common. The squat
serves a page of syndicated advertising and monetises accidental arrivals a
fraction of a cent at a time. Individually worthless; at scale, a business.
These are the domains that resolve to a handful of well-known parking
nameservers, and a scan surfaces them in clusters — `parkingcrew.net`,
`parklogic.com`, `namefind.com` and their peers show up repeatedly across
unrelated squats, which is exactly the kind of shared infrastructure that makes
[campaign clustering]({{ site.baseurl }}/internals/analysis/) work.

**Phishing.** The lookalike hosts a credential form. The name is doing the work
that a forged certificate used to do: it makes the address bar read close enough
to right. Typosquats are used here less as a landing point for typing errors and
more as a *plausible* string to put in an email, where the reader is checking
the domain against memory rather than character by character.

**Malware and drive-by delivery.** Historically a download that the user
believes comes from the real project. In the package-registry era this has
largely moved to the supply chain, below.

**Supply-chain compromise.** The squat is not a website at all; it is a package
on npm, PyPI, RubyGems, Maven, crates.io or a Go module proxy, named so that a
developer or a build machine installs it by mistake. The payload runs at install
time with the developer's credentials, or inside CI with the pipeline's. This is
the highest-consequence variant and the one that has grown fastest, because the
victim is not a person reading a URL but a machine following a manifest. It is
covered on its own in [chapter 2](../surface/).

**Brand extortion and resale.** The squat is registered to be sold to the
trademark holder, or held to embarrass them. Legal process exists for this
(UDRP, the ACPA in the US, national equivalents) and works about as well as
legal process against anything cheap, distributed and anonymous.

**Interception.** Email is the underrated case. A domain one character off yours
with a catch-all MX record receives every message a correspondent misaddresses:
invoices, contracts, password resets, calendar invitations. Nothing needs to be
clicked. This is why `dns-mx` matters more in a typosquatting scan than its
prominence in most tools suggests — a squat with MX records configured is a
qualitatively different finding from one that merely parks.

**Sabotage and traffic theft.** A competitor's name, redirected. Cruder, rarer,
and usually shorter-lived because it is easy to attribute.

## The vocabulary

The literature has accumulated a set of terms which are used loosely elsewhere
and precisely here. URLInsane's algorithm names follow these distinctions.

**Typosquatting** proper: the variant is a *typing error* away from the target.
`exmaple.com` for `example.com`. Motor error, keyboard geometry, single edits.

**Combosquatting**: the variant is the target name plus a real word —
`example-login.com`, `secure-example.com`. No typing error is involved at all,
which is what makes it nastier: there is nothing misspelled for a careful reader
to catch, and the name is often more plausible than the real one. Kintis et
al.'s longitudinal study (CCS 2017) found these vastly outnumber typo variants
and live far longer. URLInsane covers this with `cb`.

**Bitsquatting**: the variant differs from the target by a single *bit*, not a
single keystroke — `axample.com` differs from `example.com` in one bit of one
byte. The error source is not a human at all but memory corruption in the
resolving host, which Artem Dinaburg demonstrated actually happens at
internet scale (Black Hat 2011). URLInsane covers this with `bf`, and it is the
one algorithm whose output looks like nonsense until you know why it is there.

**Homograph / homoglyph attacks**: the variant is not a different string to a
human at all, only to a machine — a Cyrillic `а` where a Latin `a` should be, or
one of the many Unicode characters that render close to `l`, `1`, `I`, `0` or
`O`. Registered as punycode (`xn--...`), rendered as the original. Covered by
`hr`.

**Doppelganger domains**: the missing-dot case. `wwwexample.com`, or
`marketingexample.com` for `marketing.example.com`. A subdomain boundary that
the typist failed to type. Covered by `do` and `di`.

**Soundsquatting**: the variant sounds the same read aloud — `example` and any
homophone of a word inside it. Matters for voice interfaces and dictation, and
for second-language users who learned a name by ear. Covered by `hs`.

**Levelsquatting**: the target name appears as a *subdomain* of an
attacker-controlled domain — `example.com.login-secure.net`. Mobile browsers
that truncate the address bar are the usual target. Covered by `si` in the
generative direction.

**Brandjacking** is the umbrella term for the commercial-harm framing of all of
the above, and **cybersquatting** is the legal term for registering a name in
bad faith with respect to a trademark, which overlaps but is not the same set —
plenty of typosquats infringe nothing and plenty of cybersquats contain no typo.

## What it is not

Two things get filed under typosquatting that behave differently and are worth
separating.

It is not **DNS hijacking**. Hijacking takes the real name and points it
somewhere else, by compromising a registrar account, a nameserver, or the
resolution path. The name is correct and the answer is wrong. Typosquatting is
the opposite: the answer is correct and the name is wrong. Defences do not
transfer.

It is not **phishing** as such. Phishing is a persuasion attack that frequently
*uses* a squatted name, but the squat is infrastructure, not the technique.
Reporting the squat and reporting the phish are different remediation paths with
different response times.

## The measurement literature

Typosquatting is unusually well studied for something so simple, because it is
one of the few abuse categories you can measure exhaustively: the candidate
space is generable, and the ground truth (registered or not) is queryable.

The work that shaped the algorithms in this tool is collected in the
[bibliography](../../reference/bibliography/), with PDFs. If you read three
things, read Szurdi et al.'s long-tail study (USENIX Security 2014), Dinaburg on
bitsquatting (Black Hat 2011), and Kintis et al. on combosquatting (CCS 2017).
Between them they establish the three claims this tool is built on: the
candidate space is enumerable, error sources include ones humans do not
generate, and the highest-volume abuse involves no typing error at all.

---

Next: **[The named-entity surface](../surface/)** — why domains are only the
oldest case.
