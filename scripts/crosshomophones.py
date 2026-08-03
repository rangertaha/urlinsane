import glob, os, re, sys, unicodedata as u
from collections import defaultdict

# Marks that carry no segmental identity: stress, length, syllable breaks, ties.
STRIP = set("ˈˌːˑ.‿|‖()[]/")
def norm_ipa(s):
    s = s.replace(" ", "")
    out = []
    for ch in s:
        if ch in STRIP:
            continue
        if u.combining(ch):          # devoicing, syllabicity, nasalisation marks
            continue
        out.append(ch)
    return "".join(out)

TYPEABLE = re.compile(r"^[a-z][a-z0-9-]{2,}$")

by_ipa = defaultdict(dict)     # ipa -> {spelling: set(langs)}
for path in sorted(glob.glob(os.path.join(sys.argv[1], "*.tsv"))):
    lang = os.path.basename(path).split("_")[0]
    with open(path, encoding="utf-8") as fh:
        for line in fh:
            parts = line.rstrip("\n").split("\t")
            if len(parts) < 2:
                continue
            word, ipa = parts[0].strip().lower(), norm_ipa(parts[1])
            if not ipa or not TYPEABLE.match(word):
                continue
            by_ipa[ipa].setdefault(word, set()).add(lang)

groups = []
for ipa, spellings in by_ipa.items():
    if len(spellings) < 2:
        continue
    langs = set()
    for ls in spellings.values():
        langs |= ls
    if len(langs) < 2:                       # within-language: hs already covers it
        continue
    words = sorted(spellings)
    if len(words) > 6:                       # a phonetic sink, not a lure
        continue
    if len(set(w.replace("-", "") for w in words)) < 2:
        continue
    if min(len(w) for w in words) < 4:       # short words collide on everything
        continue
    groups.append((ipa, words, sorted(langs)))

# Two IPA keys can normalise to the same spelling set; keep one.
seen, deduped = set(), []
for ipa, words, langs in groups:
    k = tuple(words)
    if k in seen:
        continue
    seen.add(k)
    deduped.append((ipa, words, langs))
groups = deduped
groups.sort(key=lambda g: g[1])
print(f"ipa keys={len(by_ipa)}  cross-language groups={len(groups)}", file=sys.stderr)
for ipa, words, langs in groups[:12]:
    print(f"  /{ipa}/ {' '.join(words)}   [{','.join(langs)}]", file=sys.stderr)
with open(sys.argv[2], "w", encoding="utf-8") as out:
    for ipa, words, langs in groups:
        out.write(" ".join(words) + "\n")
