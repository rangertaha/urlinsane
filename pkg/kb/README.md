# kb

Keyboard layouts as a Go library. Pick a layout, read its keys, and ask what
sits next to a character or above it on the Shift level.

```go
board, _ := kb.Get("kbdus")

board.Adjacent("e")   // [w r d 3 4 s]
board.Shifted("4")    // [$]
board.Unshifted("$")  // [4]
```

The data is the layout set published by [kbdlayout.info](https://kbdlayout.info),
covering 203 layouts across 110 languages — every keyboard Windows ships,
including the historic-script ones.

## Selecting a layout

A layout is identified by the driver it ships in, and `Get` also accepts any of
the KLIDs installed against that driver, so all of these are the same board:

```go
kb.Get("kbdus")
kb.Get("KBDUS.DLL")
kb.Get("00000409")   // en-US
kb.Get("409")        // the short form Windows tools print
kb.Get("00000804")   // zh-Hans-CN, which shares the US driver
```

To search rather than address:

```go
kb.ByLanguage("de")     // German, German (IBM), Swiss German, ...
kb.ByLanguage("de-CH")  // just the Swiss ones
kb.Find("dvorak")       // by name
kb.List()               // the whole catalogue, without parsing any of it
```

Layouts are parsed on first use and cached, so `Get` is cheap to call
repeatedly and returns the same pointer each time.

## Reading keys

Every key carries the scan code of the physical switch it sits on, the virtual
key it maps to, its position, and what it types in each modifier state.

```go
k, _ := board.Key("12")   // the E key on any layout

k.Base()      // "e"
k.Shift()     // "E"
k.AltGr()     // "" on the US layout, "€" on many European ones
k.Texts()     // every distinct string the key can type

o, _ := k.Text(kb.Caps)   // Out{Text: "E"}
o.Dead                    // true for the accent keys of a European layout
```

Caps Lock inverts Shift, so only the keys that depart from that rule — the
digit row of a German keyboard, where Caps Lock does nothing — store it. The
rest is derived, which keeps the dataset to a third of its full size.

## Adjacency

Adjacency is geometric, not tabular. Scan codes describe the physical switch
rather than the letter printed on it, so one position table serves every
layout, and neighbours are found by distance between the points a finger
would strike — key centres, except on a key wide enough that the hand meets it
wherever it happens to be. The space bar is six units across, so a
centre-to-centre reading would make `b` its only neighbour instead of the whole
bottom row.

That matters because the rows of a real keyboard are staggered. The `e` key
overhangs `d` by a quarter of a key and `s` by three quarters, so a finger
slipping down from `e` is likelier to land on `d`:

```go
board.Adjacent("e")            // [w r d 3 4 s] — d ahead of s
board.AdjacentWithin("e", 1.9) // a wider net
board.AdjacentKeys(k, 1.3)     // the keys themselves, nearest first
```

A layout stored as rows of equal-width columns cannot make that distinction,
and neither can it tell an ANSI keyboard from an ISO one:

```go
uk := kb.MustGet("kbduk")
uk.Adjacent("a")              // [s q \ z w] — ISO has a key below left Shift
uk.With(kb.ANSI).Adjacent("a") // [s q z w]  — ANSI does not
```

The result follows the layout rather than assuming QWERTY, and Shift travels
with the character:

```go
kb.MustGet("kbdgr").Adjacent("z")  // [t u h 6 7 g] — QWERTZ
kb.MustGet("kbdfr").Adjacent("a")  // [z q & é]     — AZERTY
kb.MustGet("kbddv").Adjacent("e")  // [o u . q j p] — Dvorak
kb.MustGet("kbdus").Adjacent("E")  // [W R D # $ S]
```

## Scope

The dataset covers the alphanumeric block. The function row and the navigation
cluster type nothing, and the numeric keypad only duplicates digits the number
row already has — including it would put two physical keys behind `4` and make
adjacency ambiguous.

## Storage

The dataset is protocol buffers — one encoded layout per file plus an index —
compiled into the binary with `go:embed`. The schema is
[`internal/kbpb/kb.proto`](internal/kbpb/kb.proto).

That is 257 KB for all 203 layouts, against 1.5 MB for the same data as JSON.
Selecting a layout costs only the 12 KB index; the layouts themselves are
decoded on first use and cached, so a program that touches one keyboard decodes
one keyboard.

Everything that is a number is stored as one. A scan code is the byte the
keyboard reports, a virtual key is a Windows constant, and a KLID is a 32-bit
identifier — the `"10"`, `"VK_Q"` and `"00000409"` spellings the source data
uses are presentation, and across 10,446 keys they cost 77 KB. `VirtualKey` is
an enum whose values *are* the Windows codes (`VK_A = 0x41`), so protoc's
generated name tables convert both ways for free, and a virtual key the schema
has not seen fails the build instead of being silently dropped. The API keeps
the readable spelling either way.

Key positions are not stored at all. They are recomputed from the scan code and
the layout's form when a layout loads, which keeps the geometry in one table
instead of copying it into every file.

`Layout` also marshals to JSON, which spells out everything it carries. Nothing
in the package reads that back — it is there for exporting a layout or looking
at one:

```go
raw, _ := json.MarshalIndent(kb.MustGet("kbdus"), "", "  ")
```

## Regenerating

```sh
go generate ./pkg/kb
```

That runs two steps: `protoc` for the Go bindings, then `gen`, which walks the
KLID list on kbdlayout.info, follows each to the driver that serves it, and
converts the driver's XML export. Responses are cached under `pkg/kb/.cache`,
so a re-run after a parsing change costs no requests; delete the cache to fetch
afresh.

Encoding is deterministic, so a rebuild that changed nothing leaves the dataset
byte for byte identical. Since the files are opaque to a diff, every layout is
decoded and compared before it is written, and the tests check the whole
catalogue round-trips.

Regenerating needs `protoc` and `protoc-gen-go` on `PATH`; building and using
the package does not.
