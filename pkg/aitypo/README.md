# aitypo

Packages urlinsane's typo generators as **learnable tasks**: an oracle, a
corpus, splits, and exact scoring.

```go
d := aitypo.Data{Vowels: en.Vowels(), Homoglyphs: en.Homoglyphs(), Keyboard: rows}
tasks, _ := aitypo.Tasks(d).Select("co", "cs", "hr")

ex := aitypo.Emit(tasks, corpus, "datasets/domains/domain.lst")
aitypo.Assign(ex, aitypo.DefaultRatio, "")
aitypo.WriteJSONL(os.Stdout, ex)
```

```
{"task":"co","input":"google.com","expect":["gogle.com","googe.com",...],"split":"train"}
```

## Why these functions are trainable

Every function in `pkg/typo` is a **total, deterministic map from a name to a
set of names**. `CharacterOmission("google")` is exactly
`{gogle, googe, googl, goole, oogle}` and nothing else, forever.

That is unusual, and it is the whole basis of this package. A task with an exact
oracle needs no annotation, cannot disagree with itself, generates unlimited
labelled data, and can be scored exactly rather than by a similarity metric
someone picked.

## The distinction that matters: `Needs`

Not every task asks a model the same kind of question.

| `Needs` | Tasks | What success means |
|---|---|---|
| `NeedsNothing` | `co cs cr hi ho di do dhs tos bf` | The model learned a **rule** — delete one character, transpose two, insert at every gap |
| `NeedsLanguage` | `hr hs cm vs gi gr cns ons sp` | The model memorised a **table** — this script's homoglyphs, this language's misspellings |
| `NeedsKeyboard` | `acs aci rar` | The model memorised a **layout** |

Reporting one accuracy across both would average a rule a model learned with a
lookup it memorised. For the table tasks the interesting experiment is
generalisation *off* the table — homoglyphs for a script the data does not
carry.

`sp` looks like a rule and is not. `typo.SingularPluralise` goes through
`nlp.NewClient`, so mouse/mice, child/children, person/people, datum/data and
the uncountables are dictionary entries, derivable from the input by nothing.
It sat under `NeedsNothing` until that was noticed, which credited a model with
a rule it had memorised as an English dictionary.

It is also the one table not gated on a `Data` field, because its table is
compiled into the library rather than supplied — so unlike every other
language task it is *not* silenced on a non-English corpus, and
`sp("maison")` is `["maisons"]`. Read `sp` on non-English input as English
morphology applied to foreign strings.

`bf` is a rule but not a rule about characters: it models a bit flipping in a
resolver's memory, so its output is byte-level and may not be valid UTF-8.
Learnable; not evidence about plausible human typing.

## Scoring is strict on purpose

`Result.Exact()` is the headline: the prediction must be the expected set, whole.
A generator's contract is the set — omission of `google` is five names, not four
of them and something plausible — so four right and one invented has not learned
the rule. `Precision`, `Recall` and `F1` are there to say *how* it failed.

Both sides are normalised the same way (sorted, deduplicated, input dropped), so
a model cannot inflate recall by repeating itself or echoing the input.

`Summarize` errors on a result naming a task the registry does not have, rather
than labelling it. `Registry` is a map, so a miss yields the zero `Task` and its
`Needs` is `NeedsNothing` — the zero value of the enum that separates a learned
rule from a memorised table. Scoring a corpus against `Tasks(Data{})`, which
deregisters every language and keyboard task, reported a homoglyph table and a
keyboard layout as learned rules at a perfect 1.000.

`Summarize` is **macro-averaged**: every input is one vote. Micro-averaging would
weight an input by how many variants it happens to produce, and the score would
mostly measure the corpus's length distribution.

## Splits do not leak

`Assign` partitions by **input name**, not by example. One name yields an example
per task and those share almost all their structure, so putting `google` in
train for omission and in test for transposition measures memorisation of the
name rather than of the rule.

It is a stable hash of the name, not a shuffle:

- the same name lands in the same split on every machine and every run;
- **adding names never moves the ones already there**, so a corpus can grow
  without invalidating every earlier measurement.

`Leakage` asserts the property rather than trusting it — concatenating two
corpora, or splitting with two salts, breaks it silently, and the symptom is a
test score that looks unusually good.

## It wraps `pkg/typo`, it does not copy it

This package began as a copied directory, and that is worth recording. A copy of
a generator set is the exact defect this repository spent a day retiring: `acs`,
`aci` and `rar` had been fixed rune-safe inside their own plugin directories
while eight sibling generators kept slicing bytes, and **the split was the bug**.
Four minutes after the copy was taken it was already two commits behind.

There is one implementation, in `pkg/typo`. Every oracle here calls it. A task
registry that reimplemented a generator would train a model against a second
definition of the truth.

## What cannot be written

`bf` models a bit flipping in a resolver's memory, so it flips bits *inside*
multi-byte runes: `bf("münchen")` returns 16 byte sequences out of 64 that are
not text at all.

JSON strings are UTF-8 by definition and `encoding/json` **coerces invalid input
to U+FFFD without erroring**. Written anyway, those 16 come back changed, and
scoring the corpus against its own oracle then grades a model that reproduced it
*exactly* at 0.750 precision — with no error anywhere.

So `WriteJSONL` refuses, naming the task and the input — checking every example
*before* writing any, because `bufio.Writer` flushes at 4 KB and a refusal
partway through used to leave a half-encoded record on disk. It checks `Source`
too: a path need not be valid UTF-8 on Linux, and provenance that comes back
changed is a corpus lying about where it came from. Filter with
`Representable` if some examples are expected to be unwritable, or keep such a
corpus to ASCII inputs. `bf` over ASCII is fine and is not banned.

## Cost

Generation is cheap per call and large in aggregate. 300 domains over 8 tasks is
2,382 examples and 50,682 output strings; the whole 17-task set over 2,000
domains does not finish in two minutes. `Stats.Variants` — the total size of all
expectation sets — is the figure that predicts training cost, not the example
count.

## What this package is not

It trains nothing and imports no ML machinery. `internal/model` is the learner —
an HMM with Baum-Welch and dag-cbor model artifacts — and `graph.BeliefModel` is
where a learned function plugs back into the scanner. This package only builds
the data and scores the answers.
