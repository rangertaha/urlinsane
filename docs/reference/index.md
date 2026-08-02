---
title: Reference
nav_order: 5
has_children: true
permalink: /reference/
---

# Reference

Flat material — look things up here rather than reading it through.

| Page | What it holds |
|---|---|
| [CLI]({{ site.baseurl }}/CLI/) | Every command, flag, filter, format and exit code, including the ones specified but not yet built (§9) |
| [Design]({{ site.baseurl }}/DESIGN/) | The graph engine design document in full: data model, registry, operators, plan, execution, termination |
| [Keyboards]({{ site.baseurl }}/KB/) | The `pkg/kb` layout model — geometric adjacency, layout selection, wrong-layout typing |
| [Glossary](glossary/) | Terms used across the book, defined once |
| [Bibliography](bibliography/) | The measurement literature the algorithms come from, with the PDFs |

Two Go packages are documented in place rather than here, because they are
usable on their own:

- [`pkg/kb`](https://pkg.go.dev/github.com/rangertaha/urlinsane/pkg/kb) —
  keyboard layouts with geometric key adjacency, built from
  [kbdlayout.info](http://kbdlayout.info/).
- [`pkg/typo`](https://pkg.go.dev/github.com/rangertaha/urlinsane/pkg/typo) —
  the string-level typo primitives the algorithms are built from.
