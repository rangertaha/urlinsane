---
title: 7 · Targets and scope
parent: Part II · Using URLInsane
nav_order: 3
---

# Targets and scope
{: .no_toc }

- TOC
{:toc}

## The kind is read from the string

```
urlinsane typo [<scope>] <target>
```

There is no `--type` flag. What a target *is* comes from the string alone,
because a flag that contradicted the string would only be a second way to be
wrong:

```bash
urlinsane typo acme.com                 # domain
urlinsane typo bob@acme.com             # email
urlinsane typo npm:lodash               # package
urlinsane typo github.com/acme/tool     # repo
urlinsane typo bobsmith                 # username
```

The rules are ordered most-specific first, and each is a positive test rather
than a fallthrough:

1. an `@` with a host after it → **email**
2. `host/owner/name`, or any URL or scp-style remote of one → **repo**
3. `registry:name` → **package**
4. a dotted name whose suffix is in the public suffix list → **domain**
5. anything else that is a legal handle → **username**

Rule 4 is the one that earns its keep. `lodash` and `lodash.com` are both legal
hostnames, so a hostname test alone would classify every bare package name as a
domain and go do DNS lookups for it. Requiring a real registry suffix is the
only discriminator that does not need a list of known package names — which
means a bare `lodash` is read as a **username**, and you must write `npm:lodash`
to get a package.

Rule 3 tests the registry token structurally rather than against a list: a
registry name is a bare word, so a dot or a slash before the colon means the
colon is doing some other job. `example.com:8080` is therefore not package
`8080` on registry `example.com`, and `https://x` is not a package on registry
`https`. A closed list of known registries would have been worse — it would
reject every registry added after the line was written.

### URLs are stripped

A pasted URL is the common case for a domain target, so the scheme, path, port
and query are removed before the hostname test:

```bash
urlinsane typo https://example.com/about?utm=x    # same as: urlinsane typo example.com
```

They carry no identity. `https://example.com/about` and `example.com` are the
same name, and minting separate nodes for them would defeat the convergence the
whole graph model rests on.

### What is refused

```console
$ urlinsane typo 192.0.2.1
scan: "192.0.2.1" is an IP address; there is nothing to typosquat — scan the name that resolves to it
```

An address is a real entity but not a *nameable* one. Varying `192.0.2.1`
character by character produces addresses with no relationship to the original.
The type system enforces this: `ip` is an **observed** type, not a nameable one,
so it can never root a variant. See [Types and relations](../../internals/types/).

If nothing matches at all:

```console
$ urlinsane typo "not a name!!"
scan: cannot tell what "not a name!!" is: expected an email, a repo URL, a registry:package, a domain or a username
```

## Decomposition: one target, several names

A target is rarely one name. Before any variant is generated, the seed is
decomposed into the entities it is made of:

| Target | Decomposes into |
|---|---|
| `acme.com` | the domain, and its TLD |
| `bob@acme.com` | the address, the local part `bob` (a username), and `acme.com` (a domain) |
| `npm:lodash` | the package, and its owner |
| `github.com/acme/tool` | the repo, the host it is on, and the owner `acme` |

Then **every nameable node in that closure is varied** — including the composite
seed itself. So a bare email target yields whole-address variants
(`bob@acme-corp.com`) *as well as* local-part variants (`bobb@acme.com`) and
domain variants (`bob@acmecom.com` and friends). Three attacks, one scan; see
[the named-entity surface](../../attack/surface/) for why they are different
attacks.

Structural edges — `TLD_OF`, `LOCAL_PART`, `DOMAIN_OF`, `OWNER`, `HOSTED_ON` —
cost **zero depth**, because they are derived from parsing the target string and
required no network call. Decomposition therefore never eats into your `--depth`
budget.

## Scope: narrowing what gets varied

The optional first positional narrows the set of nodes that get varied. It does
not change how the target is read:

```bash
urlinsane typo username bob@acme.com        # vary only bob
urlinsane typo domain   bob@acme.com        # vary only acme.com
urlinsane typo username,domain bob@acme.com # both, but not the whole address
```

This matters because the parse and the scope are deliberately independent.
`urlinsane typo username bob@acme.com` and `urlinsane typo bob@acme.com` parse
the target *identically* — an email is an email either way. If scope changed the
parse, the same target string would mean different things on different runs and
the seed closure would differ with it, which would make two scans of "the same
thing" incomparable.

Scope accepts any nameable type name: `domain`, `email`, `package`, `repo`,
`username`. Comma-separate for several.

### When to use it

- **Auditing a person, not a company.** `urlinsane typo username bob@acme.com`
  keeps the scan on `bob` and off the corporate domain you have already scanned.
- **Auditing the domain of an address you were given.**
  `urlinsane typo domain bob@acme.com` — no need to strip the string by hand.
- **Cutting the cost.** Varying a three-part repo target produces variants of
  the host, the owner and the name. If you only care about the owner handle,
  say so.

## Choosing a good seed

Two practical notes.

**Scan the registrable domain, not a subdomain.** `www.acme.com` and `acme.com`
generate different neighbourhoods, and the second is the one attackers register.
Subdomain-shaped attacks are covered by the `si` algorithm from the registrable
name, not by seeding the subdomain.

**Scan each brand name separately.** If your product is `acme` and your domain
is `acmecorp.com`, the neighbourhood of `acme` on package registries has nothing
to do with the neighbourhood of `acmecorp.com` in DNS. They are different seeds.

---

Next: **[Algorithms](../algorithms/)**.
