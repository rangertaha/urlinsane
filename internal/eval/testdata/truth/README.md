# Truth sets

JSONL, one `eval.Record` per line. Each record is a name that *exists* and
imitates a brand.

## `example.jsonl` is a fixture, not a measurement

It is hand-written to exercise the scorer — a few names each algorithm family
should reach, plus some nothing reaches. **Do not quote a recall number from
it.** Hand-authored squats are drawn from what the algorithms already generate,
which is exactly the bias recall is supposed to detect.

## Building a real truth set

```sh
go run ./cmd/eval fetch --brand paypal.com \
    --exclude paypal-corp.com --exclude paypalobjects.com \
    --out internal/eval/testdata/truth/paypal.jsonl
```

Then **read the file**. Fetched records are written with `"reviewed": false`
because Certificate Transparency answers "what names exist", not "what names are
squatting": a substring query returns the brand's own vanity domains and
unrelated names that share a substring. Set `"reviewed": true` on the ones you
confirm and delete the rest.

crt.sh is frequently overloaded and answers 502/503. A failed fetch means "try
later", not "no squats exist" — the distinction matters, because an empty truth
set would otherwise score as perfect recall. `fetch` errors out on an empty
result for that reason.

## Scoring

```sh
go run ./cmd/eval score --truth internal/eval/testdata/truth --missed
```

The primary metric is **core recall**: the fraction of truth names whose
registrable label the algorithms generate. It deliberately ignores the public
suffix, because reaching `paypal-login.net` from `paypal.com` needs a
combosquat *composed with* a TLD swap, and `score` applies each algorithm once.
A real scan composes across rounds, so this number is a lower bound.

The per-algorithm table is the actionable half. An algorithm with many
candidates and zero unique hits is paying for coverage something else already
provides.
