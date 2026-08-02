---
title: Part I · The attack
nav_order: 2
has_children: true
permalink: /attack/
---

# Part I · The attack

Nothing in this part is about URLInsane. It is about the thing URLInsane looks
for, which existed long before the tool and would carry on without it.

Typosquatting sits in an awkward place. It is not a software vulnerability —
there is no bug to patch, no version to upgrade past, no CVE. It is not really
social engineering either, because in the ordinary case nobody is persuaded of
anything; the victim does exactly what they intended to do, correctly, at the
wrong name. It is closest to a *systems* problem: human naming is
approximate, machine name resolution is exact, and an attacker lives in the gap.

Four chapters:

1. **[What typosquatting is](typosquatting/)** — the mechanism, the money, and
   the taxonomy of what people do with a stolen name.
2. **[The named-entity surface](surface/)** — why domains are only the oldest
   case, and what changes when the namespace is npm, PyPI, GitHub, or a company
   directory.
3. **[Why names get mistyped](errors/)** — motor slips, perceptual confusions,
   knowledge errors, and hardware faults. This chapter is the reason the
   algorithm list looks the way it does.
4. **[What defenders do about it](defending/)** — the point of a variant list,
   and the honest limits of one.
