# kb

Keyboard layouts as a Go library. Pick a layout, read its keys, ask what sits
next to a character or above it on the Shift level, and work out what someone
typed with the wrong layout selected.

```go
board := kb.MustGet("kbdus")

board.Adjacent("e")   // [w r d 3 4 s]
board.Shifted("4")    // [$]
board.Unshifted("$")  // [4]
```

The data is the layout set published by [kbdlayout.info](https://kbdlayout.info):
**203 layouts across 110 languages** — every keyboard Windows ships, including
the historic-script ones.

Every snippet below is taken from a runnable example in `example_test.go`, so
the outputs are checked by `go test` rather than written from memory.

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
kb.ByLanguage("de")      // German, German (IBM), German Extended ×2, Swiss German
kb.ByLanguage("de-CH")   // just the Swiss one
kb.Find("dvorak")        // by name
kb.List()                // the whole catalogue, without decoding any of it
kb.IDs()                 // just the driver names
kb.Languages()           // the 110 subtags, sorted
```

A bare subtag matches every layout for the language; a full BCP-47 tag matches
only the layouts carrying that exact locale. All of these return `nil` when
nothing matches.

You can also ask which keyboards can actually type something:

```go
kb.ByString("google.com")   // 98 layouts
kb.ByString("münchen")      // 20
kb.ByString("中文")          // 0
kb.ByKeys("a", "b", "c")    // the same question, one text per argument
```

Both intersect — a layout has to type *every* text given. Unlike the lookups
above these cannot be answered from the index, so the first call decodes the
whole catalogue (about 50 ms); afterwards it is a map lookup per layout.

Layouts are decoded on first use and cached, so `Get` is cheap to call
repeatedly and returns the same pointer each time. Nothing it returns may be
modified.

## Reading keys

Every key carries the scan code of the physical switch it sits on, the virtual
key it maps to, its position, and what it types in each modifier state.

```go
k, _ := kb.MustGet("kbdgr").Key("29")

k.SC     // "29"
k.VK     // "VK_OEM_5"
k.Base() // "^"
k.Shift()// "°"
k.AltGr()// "" here; "€" on many European layouts
k.Texts()// every distinct string the key can type

o, _ := k.Text(kb.Base)
o.Dead   // true — the German circumflex arms an accent, it does not type one
```

Caps Lock inverts Shift, so only the keys that depart from that rule — the
digit row of a German keyboard, where Caps Lock does nothing — store it. The
rest is derived.

## Adjacency

Adjacency is geometric, not tabular. Scan codes describe the physical switch
rather than the letter printed on it, so one position table serves every
layout, and neighbours are found by distance between the points a finger would
strike.

That matters because the rows of a real keyboard are staggered. The `e` key
overhangs `d` by a quarter of a key and `s` by three quarters, so a finger
slipping down from `e` is likelier to land on `d`, and the result is ordered
accordingly:

```go
us.Adjacent("e")             // [w r d 3 4 s]  — d ahead of s
us.Adjacent("E")             // [W R D # $ S]  — Shift travels with it
us.AdjacentWithin("e", 1.9)  // a wider net
us.AdjacentKeys(k, 1.3)      // the keys themselves, nearest first
```

A table of equal-width rows cannot make that distinction. Neither can it tell
an ANSI keyboard from an ISO one, where the short left Shift makes room for a
key that neighbours `a`:

```go
uk := kb.MustGet("kbduk")
uk.Adjacent("a")               // [s q \ z w]  — ISO
uk.With(kb.ANSI).Adjacent("a") // [s q z w]    — ANSI has no such key
```

Nor can it handle a key that is not one unit wide. The space bar is six units
across, so measuring to its centre would make `b` its only neighbour instead of
the whole bottom row; distance is measured to the nearest point a hand would
actually strike, which is the centre for every ordinary key.

And the answer follows the layout rather than assuming QWERTY:

```go
kb.MustGet("kbdgr").Adjacent("z")  // [t u h 6 7 g]  — QWERTZ
kb.MustGet("kbdfr").Adjacent("a")  // [z q & é]      — AZERTY
kb.MustGet("kbddv").Adjacent("e")  // [o u . q j p]  — Dvorak
```

## Typing on the wrong layout

A `Stroke` is a key and the modifiers held with it. It says nothing about what
gets typed — that depends on the layout the keystroke lands on, which is the
whole point:

```go
us.Strokes("abc")            // [{1E base} {30 base} {2E base}]
ru.Type(us.Strokes("abc"))   // "фис"
```

`Translate` does both steps at once:

```go
us.Translate("hello", ru)       // "руддщ"
us.Translate("google.com", ru)  // "пщщпдуюсщь"
us.Translate("google.com", fr)  // "google:co," — AZERTY moves the punctuation
```

Characters the source layout cannot type, and keys the target leaves bare, are
passed through unchanged rather than dropped, so a domain keeps its dots. Two
consequences worth knowing: the result is **not always the same length**, since
a few layouts have ligature keys that type two characters at once; and it
**round-trips only if the source can type everything given**, because a
character that falls through untouched is an ordinary character on the way
back.

## Printing

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

Boxes are placed from the real geometry, so the indentation between rows is the
physical stagger rather than decoration. Shift sits above base as on a keycap,
and dead keys are bracketed because they type nothing on their own.

`Print` writes to stdout, `Fprint` to any writer, and `String` returns the same
drawing — so `fmt.Print(layout)` agrees with all of them.

## Scope and known gaps

The dataset covers the **alphanumeric block**. The function row and the
navigation cluster type nothing, and the numeric keypad only duplicates digits
the number row already has — including it would put two physical keys behind
`4` and make adjacency ambiguous.

Three things are missing, and they are worth knowing before trusting a negative
answer:

- **Dead-key compositions are not stored.** The dataset records that a key is
  dead and what accent it carries, but not what it composes into. So German
  counts as unable to type `é`, and `ByString("café")` leaves it out, even
  though `´` then `e` produces it. `Translate` likewise emits the bare accent
  where a real keyboard would arm one.
- **`kbdjpn` is Latin-only here.** Its kana sit behind a Kana lock, a modifier
  this package has no bit for, so 97 of its keystrokes are absent.
- **`kbdcan` is missing its `VK_OEM_8` level** — the ¹ ² ³ ¼ row, 78
  keystrokes, for the same reason.

Those are the only two layouts affected, and rebuilding the dataset reports any
level it cannot represent, so a new one cannot slip in unnoticed.

## Storage

The dataset is protocol buffers — one encoded layout per file plus an index —
compiled into the binary with `go:embed`. The schema is
[`internal/kbpb/kb.proto`](internal/kbpb/kb.proto).

That is 257 KB for all 203 layouts, against 1.5 MB for the same data as JSON.
Selecting a layout costs only the 12 KB index; the layouts themselves are
decoded on first use, so a program that touches one keyboard decodes one
keyboard.

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
afresh. Cache entries are written atomically, so an interrupted run cannot
leave a truncated one behind to be reused.

Encoding is deterministic, so a rebuild that changed nothing leaves the dataset
byte for byte identical. Since the files are opaque to a diff, every layout is
decoded and compared before it is written, and the tests check the whole
catalogue round-trips.

Regenerating needs `protoc` and `protoc-gen-go` on `PATH`; building and using
the package does not.
