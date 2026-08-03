---
title: Keyboards
parent: Reference
nav_order: 3
---

# Keyboards

`pkg/kb` is URLInsane's keyboard layout library. It answers three questions
that typosquatting depends on:

- **What sits next to this character?** — the fat-finger substitutions and
  insertions a typist actually makes.
- **Which keyboards can type this at all?** — whether a variant is reachable
  by the people you are modelling.
- **What does this become on the wrong layout?** — what someone types when the
  keyboard in the OS is not the one under their hands.

It covers **203 layouts across 110 languages**, taken from
[kbdlayout.info](https://kbdlayout.info), which publishes the key tables of
every keyboard driver Windows ships.

Table of contents
=================

* [Why not a grid of rows](#why-not-a-grid-of-rows)
* [How adjacency is measured](#how-adjacency-is-measured)
* [Choosing a keyboard](#choosing-a-keyboard)
* [Wrong-layout typing](#wrong-layout-typing)
* [Seeing a layout](#seeing-a-layout)
* [What the data holds](#what-the-data-holds)
* [What is missing](#what-is-missing)
* [Rebuilding the dataset](#rebuilding-the-dataset)

## Why not a grid of rows

The obvious way to store a keyboard is as rows of characters:

```
1234567890-
qwertyuiop
asdfghjkl
zxcvbnm
```

and to call two characters adjacent when they are next to each other in that
grid. It is how the earlier implementation worked, and it is wrong in three
ways that matter.

**Rows are staggered.** On a real keyboard the `qwerty` row sits half a key to
the right of the number row, and the `asdf` row a further quarter. So the key
below `e` is not the one in the same column. `e` overhangs `d` by a quarter of
a key and `s` by three quarters, which makes `d` the likelier slip — and a grid
of equal-width columns cannot express the difference at all.

**Keyboards are not all the same shape.** A European ISO keyboard has a key
between the left Shift and `z` that a US ANSI keyboard does not, and it
neighbours `a`. In a row-of-strings model that key has nowhere to live.

**Not every key is one unit wide.** The space bar is six units across and runs
the length of the bottom row. Anything that reduces it to a single grid cell
gets its neighbours wrong.

## How adjacency is measured

Every key carries the **scan code** of the physical switch it sits on. A scan
code is a property of the keyboard, not of the letter printed on the key: `10`
is the switch one row up and one column right of Caps Lock whether it types
`q`, `a`, or `ф`. So one position table serves all 203 layouts, and each layout
only has to say what its keys produce.

Positions are measured in **key units** — one unit is the width of a plain
letter key — and adjacency is the distance between the points a finger would
strike, with anything inside 1.3 units counting as a neighbour:

```go
us := kb.MustGet("kbdus")
us.Adjacent("e")   // [w r d 3 4 s]
```

The order is by distance: `w` and `r` at 1.0, then `d` at 1.031, then `3` and
`4` at 1.118, then `s` at 1.25. That `d` outranks `s` is the stagger showing
through.

Shift travels with the character, so the neighbours of a capital are capitals:

```go
us.Adjacent("E")   // [W R D # $ S]
```

The answer follows the layout rather than assuming QWERTY:

```go
kb.MustGet("kbdgr").Adjacent("z")  // [t u h 6 7 g]  — QWERTZ
kb.MustGet("kbdfr").Adjacent("a")  // [z q & é]      — AZERTY
kb.MustGet("kbddv").Adjacent("e")  // [o u . q j p]  — Dvorak
```

Form changes the answer where the boards differ. The UK layout has 52 keys on
ISO and 51 on ANSI, and the extra one is a neighbour of `a`:

```go
uk := kb.MustGet("kbduk")
uk.Adjacent("a")                // [s q \ z w]
uk.With(kb.ANSI).Adjacent("a")  // [s q z w]
```

Width is handled the same way. A key wider than one unit is struck wherever
the hand already is, so the space bar reaches the whole bottom row rather than
just the key above its centre.

## Choosing a keyboard

A layout is named after the Windows driver it ships in, and also answers to any
of the KLIDs installed against that driver — one driver often serves several
locales, so `kbdus` backs English, Bulgarian and three Chinese locales alike:

```go
kb.Get("kbdus")       // by driver
kb.Get("KBDUS.DLL")   // by file
kb.Get("00000409")    // by KLID
kb.Get("409")         // the short form Windows tools print
```

To search rather than address:

```go
kb.ByLanguage("de")       // every German layout — there are five
kb.ByLanguage("de-CH")    // only the Swiss one
kb.Find("dvorak")         // by name
kb.List()                 // the catalogue, without decoding any of it
kb.Languages()            // the 110 subtags
```

Or by what a keyboard has to be able to type — useful for deciding whether a
generated variant is reachable at all:

```go
kb.ByString("google.com")   // 98 layouts
kb.ByString("münchen")      // 20
kb.ByString("中文")          // 0
```

## Wrong-layout typing

A `Stroke` is a key and the modifiers held with it. It carries no text of its
own — what it types depends on the layout it lands on, which is the point:

```go
us, ru := kb.MustGet("kbdus"), kb.MustGet("kbdru")

us.Strokes("abc")            // [{1E base} {30 base} {2E base}]
ru.Type(us.Strokes("abc"))   // "фис"
```

`Translate` does both steps together. This is a squatting vector in its own
right — a familiar name typed with the wrong keyboard selected produces a
different, registrable string:

```go
us.Translate("hello", ru)       // "руддщ"
us.Translate("google.com", ru)  // "пщщпдуюсщь"
us.Translate("google.com", fr)  // "google:co,"
```

The French case is the one to notice: AZERTY moves the punctuation, so
`google.com` becomes `google:co,` without a single letter changing. US to
German leaves `google.com` alone, because QWERTZ only swaps `y` and `z` — the
correct null result.

Characters the source cannot type pass through unchanged, so a domain keeps its
dots. Two consequences follow: the result is not always the same length, since
a few layouts have ligature keys that type two characters at once; and it
round-trips only when the source layout can type every character it was given.

## Seeing a layout

```go
kb.MustGet("kbdgr").Print()
```

```
kbdgr — German (iso)
+---+---+---+---+---+---+---+---+---+---+---+---+---+
| ° | ! | " | § | $ | % | & | / | ( | ) | = | ? |[`]|
|[^]| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | ß |[´]|
+-----+---+---+---+---+---+---+---+---+---+---+---+---+
|     | Q | W | E | R | T | Z | U | I | O | P | Ü | * |
|     | q | w | e | r | t | z | u | i | o | p | ü | + |
+------+---+---+---+---+---+---+---+---+---+---+---+---+
|      | A | S | D | F | G | H | J | K | L | Ö | Ä | ' |
|      | a | s | d | f | g | h | j | k | l | ö | ä | # |
+----+---+---+---+---+---+---+---+---+---+---+---+-+---+
|    | > | Y | X | C | V | B | N | M | ; | : | _ |
|    | < | y | x | c | v | b | n | m | , | . | - |
+----+---+---+-+------------------------++---+---+
               |                        |
               |                        |
               +------------------------+
```

The boxes are placed from the real geometry, so the indentation between rows is
the physical stagger rather than decoration — it is the quickest way to check
that a layout loaded the way you expected. Shift sits above base as on a
keycap, and dead keys are bracketed because they type nothing on their own.

## What the data holds

| | |
|---|---|
| Layouts | 203 |
| Languages | 110 |
| Locales (KLIDs) | 218 |
| Keys | 10,446 |
| Form factors | 110 ANSI, 93 ISO |
| Size on disk | 257 KB |
| Catalogue index | 12 KB |

The dataset is protocol buffers — one encoded layout per file plus an index —
compiled into the binary with `go:embed`. As JSON the same data was 1.5 MB.
Selecting a layout costs only the index; layouts are decoded on first use and
cached, so a program that touches one keyboard decodes one keyboard.

Everything that is already a number is stored as one. A scan code is the byte
the keyboard reports, a virtual key is a Windows constant, and a KLID is a
32-bit identifier; the `"10"`, `"VK_Q"` and `"00000409"` spellings the source
uses are presentation, and across 10,446 keys they cost 77 KB. Key positions
are not stored at all — they are recomputed from the scan code and the form
when a layout loads, which keeps the geometry in one table rather than copied
into every file.

Only the **alphanumeric block** is covered. The function row and the navigation
cluster type nothing, and the numeric keypad only duplicates digits the number
row already has — including it would put two physical keys behind `4` and make
adjacency ambiguous.

## What is missing

Three gaps are worth knowing before trusting a negative answer.

**Dead-key compositions are not stored.** The data records that a key is dead
and which accent it carries, but not what it composes into. So German counts as
unable to type `é`, and `ByString("café")` leaves it out, even though `´` then
`e` produces it. `Translate` likewise emits the bare accent where a real
keyboard would arm one.

**`kbdjpn` is Latin-only here.** Its kana sit behind a Kana lock, a modifier
the library has no bit for, so 97 of its keystrokes are absent.

**`kbdcan` is missing its `VK_OEM_8` level** — the ¹ ² ³ ¼ row, 78 keystrokes,
for the same reason.

Those are the only two layouts affected, and rebuilding reports any modifier
level it cannot represent, so a new one cannot slip in unnoticed.

## Rebuilding the dataset

```sh
go generate ./pkg/kb
```

Two steps run: `protoc` for the Go bindings, then the generator, which walks
the KLID list on kbdlayout.info, follows each to the driver that serves it, and
converts that driver's XML export. Responses are cached under `pkg/kb/.cache`,
so a re-run after a parsing change costs no requests; delete the cache to fetch
afresh.

Encoding is deterministic, so a rebuild that changed nothing leaves the dataset
byte for byte identical. Because the files are opaque to a diff, every layout
is decoded and compared before it is written, and the tests check that the
whole catalogue round-trips.

Rebuilding needs `protoc` and `protoc-gen-go` on `PATH`. Building and using the
library does not.

---

Full API reference: `go doc github.com/rangertaha/urlinsane/pkg/kb`, or the
package's own [README](https://github.com/rangertaha/urlinsane/tree/master/pkg/kb).
