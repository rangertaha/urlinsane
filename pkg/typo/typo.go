// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package typo

import (
	"sort"
	"strings"
	"sync"

	"github.com/rangertaha/urlinsane/pkg/nlp"
)

// CharacterSwapping refers to a type of typo where two adjacent characters in
// the original token are exchanged or swapped. This often occurs when characters
// are unintentionally reversed in order, resulting in a misspelling.For example,
// the word "example" could become "examlpe" by swapping the position of the
// letters "l" and "p".
func CharacterSwapping(token string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i := 0; i+1 < len(rs); i++ {
		u.add(token, joinRunes(rs[:i], rs[i+1], rs[i], rs[i+2:]))
	}
	return u.tokens()
}

// AdjacentCharacterSubstitution typos happen when a character in the original
// token is mistakenly replaced by a neighboring character from the same keyboard
// layout. This type of error often occurs due to hitting an adjacent key by accident.
// For example, the token "ezample" contains a typo where the letter "x" is
// substituted with "z," which is the neighboring key on an English QWERTY
// keyboard layout.
func AdjacentCharacterSubstitution(token string, keyboard ...string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i, r := range rs {
		for _, key := range adjacentCharacters(string(r), keyboard...) {
			u.add(token, joinRunes(rs[:i], key, rs[i+1:]))
		}
	}
	return u.tokens()
}

// AdjacentCharacterInsertion typos occur when characters adjacent of each
// letter are inserted. For example, googhle inserts "h" next to it's
// adjacent character "g" on an English QWERTY keyboard layout.
func AdjacentCharacterInsertion(token string, keyboard ...string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i, r := range rs {
		for _, key := range adjacentCharacters(string(r), keyboard...) {
			u.add(token, joinRunes(rs[:i], key, r, rs[i+1:]))
			u.add(token, joinRunes(rs[:i], r, key, rs[i+1:]))
		}
	}
	return u.tokens()
}

// HyphenInsertion typos happen when hyphens are mistakenly placed between
// characters in a token, often occurring in various positions around the
// letters. This type of error can lead to unnecessary fragmentation of the
// word, with hyphens inserted at different points throughout the token.
// For example, the word "example" might be incorrectly written as "-example",
//
//	"e-xample", "ex-ample", "exa-mple", "exam-ple", "examp-le", or even
//
// "example-", with hyphens appearing before, between, or after the letters.
func HyphenInsertion(token string) (tokens []string) {
	return insertEverywhere(token, "-")
}

// HyphenOmission typos occur when hyphens are unintentionally left out of a
// token, resulting in a version of the token that misses the expected hyphenation.
// For example, the token "one-for-all" might be mistakenly written as "onefor-all",
// "one-forall", or even "oneforall", where the hyphens are omitted.
func HyphenOmission(token string) (tokens []string) {
	return characterDeletion(token, "-")
}

// DotInsertion typos take place when periods (.) are mistakenly added at
// various points within a token, leading to an incorrect placement of dots in
// the target token. This type of error typically happens due to inadvertent
// key presses or misplacement while typing. For instance, the word "example"
// may be mistakenly written as "e.xample", "ex.ample", "exa.mple", "exam.ple",
// or "examp.le", where the dot appears at different locations
// within the token, disrupting the original structure.
func DotInsertion(token string) (tokens []string) {
	// A leading or trailing dot is not a name, so those two positions are
	// dropped rather than trimmed back onto the token: trimming turned them
	// into the token itself, and DotInsertion("ab") returned nothing but "ab".
	u := newUniq()
	for _, v := range insertEverywhere(token, ".") {
		if strings.HasPrefix(v, ".") || strings.HasSuffix(v, ".") {
			continue
		}
		u.add(token, v)
	}
	return u.tokens()
}

// DotOmission typos happen when periods (.) that should be present in the target
// token are unintentionally omitted or left out. This type of error typically
// occurs when the user fails to input the expected dots, often resulting in a
// word or sequence that appears as a single string without proper separation.
// For example, the sequence "one.two.three" might be mistakenly written
// as "one.twothree", "onetwo.three", or even "onetwothree", where the dots
// are missing between certain parts of the token, causing it to lose the
// intended structure or meaning.
func DotOmission(token string) (tokens []string) {
	return characterDeletion(token, ".")
}

// GraphemeInsertion, also known as alphabet insertion, occurs when one or more
// unintended letters are added to a valid token, leading to a modified or
// misspelled version of the original token. These extra characters are typically
// inserted either at the beginning or within the token, causing it to deviate
// from its intended form. This type of error is often the result of a slip
// of the finger or an accidental keystroke. For example, the token "example"
// might be mistakenly typed as "aexample", "eaxample", "exaample", "examaple",
//
//	or "eaxampale", where additional letter like "a" are inserted throughout
//
// the token, distorting its original structure.
func GraphemeInsertion(token string, graphemes ...string) (tokens []string) {
	u := newUniq()
	for _, g := range graphemes {
		for _, v := range insertEverywhere(token, g) {
			u.add(token, v)
		}
	}
	return u.tokens()
}

// GraphemeReplacement, also known as alphabet replacement, occurs when characters
// from the original token are replaced by other letters from the alphabet,
// resulting in a modified version of the token. This type of error typically leads
// to small changes in the original token, where one or more letters are swapped
// for different characters. For example, the token "example" could be mistakenly
// written as "axample", "bxample", "cxample", "dxample", or "eaample", where
// letters like "a", "b", "c", "d", or "e" are substituted, altering the
// word slightly but keeping its general structure.
func GraphemeReplacement(token string, graphemes ...string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i := range rs {
		for _, g := range graphemes {
			u.add(token, joinRunes(rs[:i], g, rs[i+1:]))
		}
	}
	return u.tokens()
}

// CharacterRepetition typos occur when a letter is unintentionally repeated
// within a token, leading to a misspelled version. This type of error typically
// happens when a key is pressed twice or a letter is accidentally duplicated.
// For example, the token "example" might be mistakenly written as "eexample",
// "exaample", "exammple", "examplee", or "examplle", where one or more
// characters are repeated, causing the token to diverge from its original form.
func CharacterRepetition(token string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i := range rs {
		u.add(token, joinRunes(rs[:i], rs[i], rs[i], rs[i+1:]))
	}
	return u.tokens()
}

// RepetitionAdjacentReplacement typos occur when consecutive, identical letters
// in a token are replaced with adjacent keys on the keyboard, resulting in a
// slight alteration of the original word. This type of error often happens due
// to accidental key presses of nearby characters. For example, the token
// "google" might be mistakenly typed as "gppgle" or "giigle", where the repeated
// letters are swapped with neighboring keys on the keyboard, causing the word
// to be misspelled.
func RepetitionAdjacentReplacement(token string, keyboard ...string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i := 0; i+1 < len(rs); i++ {
		if rs[i] != rs[i+1] {
			continue
		}
		for _, key := range adjacentCharacters(string(rs[i]), keyboard...) {
			u.add(token, joinRunes(rs[:i], key, key, rs[i+2:]))
		}
	}
	return u.tokens()
}

// CharacterOmission occurs when one character is unintentionally omitted from
// the token, leading to an incomplete version of the original word. This type
// of typo can happen when a key is accidentally skipped or overlooked while
// typing. For example, the word "google" might be mistakenly written as "gogle",
// "gogle", "googe", "googl", "goole", or "oogle", where a single character is
// missing from different positions in the word, causing it to deviate from
// the correct spelling.
func CharacterOmission(token string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i := range rs {
		u.add(token, joinRunes(rs[:i], rs[i+1:]))
	}
	return u.tokens()
}

// pluralizer is the inflection client, built once and shared by every call.
//
// nlp.NewClient loads every irregular, plural, singular and uncountable rule
// and compiles a regexp, and SingularPluralise built a fresh one on every
// invocation: 500us per name against 4us for the next slowest generator in this
// package, 113x, to rebuild a value that is identical every time. sp runs once
// per candidate, so a scan paid for the whole rule set once per name it
// generated.
//
// Sharing it is safe rather than merely convenient: the Add*Rule methods are
// the only ones that write to a Client, nothing here calls them, and
// Plural/Singular/IsPlural/IsSingular only read. OnceValue makes the build
// happen once even if the engine runs generators concurrently.
var pluralizer = sync.OnceValue(nlp.NewClient)

// SingularPluralise typos are where a word is altered by switching between its
// singular and plural forms. This subtle change can create a word that looks
// similar to the original, but with a small variation that is easy to overlook.
// For example, if the original word is 'example', a Singular-Plural might result
// in 'examples', or vice versa.
func SingularPluralise(token string) (tokens []string) {
	// The client reports "" as both singular and plural and inflects it to ""
	// twice, which is two empty variants of an empty name.
	if token == "" {
		return nil
	}
	// Through uniq like every other generator here, rather than appending
	// straight to the result. An uncountable noun is reported as both singular
	// and plural and inflects to itself either way, so "mail", "news", "media"
	// and "famous" each came back as two copies of the name itself — a variant
	// equal to its origin, twice.
	pluralize := pluralizer()
	u := newUniq()
	if pluralize.IsPlural(token) {
		u.add(token, pluralize.Singular(token))
	}
	if pluralize.IsSingular(token) {
		u.add(token, pluralize.Plural(token))
	}
	return u.tokens()
}

// CommonMisspellings refers to typos created by frequent spelling errors or
// missteps that occur in the target language. These errors often involve slight
// changes to the spelling of a word, making them appear similar to the original
// but incorrect. For instance, the word "youtube" could be mistyped as
// "youtub", and "abseil" could become "absail", where common mistakes in
// spelling lead to slightly altered but recognizable versions of the original.
func CommonMisspellings(token string, dataset ...[]string) (words []string) {
	return swapWordSets(token, dataset)
}

// swapWordSets substitutes each member of a word set for every other member it
// finds in the token. Sets often contain members that are substrings of each
// other ("hola"/"ola", "adz"/"adze"); substituting the shorter one inside a
// match of the longer one yields nonsense ("hola" -> "hhola"), so a member is
// skipped when a longer member of the same set also matches the token.
func swapWordSets(token string, dataset [][]string) (words []string) {
	// Deduplicated like every other generator here. Two different sets can
	// contain the same pair — the cross-language homophone data stores one
	// group per pronunciation, and a word with two pronunciations appears in
	// both — so the same substitution is reachable twice and the caller would
	// get the same variant twice.
	u := newUniq()
	for _, wordset := range dataset {
		for _, word := range wordset {
			if !strings.Contains(token, word) || shadowed(token, word, wordset) {
				continue
			}
			for _, w := range wordset {
				if w != word {
					u.add(token, strings.Replace(token, word, w, -1))
				}
			}
		}
	}
	return u.tokens()
}

// shadowed reports whether word is a proper substring of a longer member of
// wordset that also occurs in token.
func shadowed(token, word string, wordset []string) bool {
	for _, other := range wordset {
		if len(other) > len(word) && strings.Contains(other, word) && strings.Contains(token, other) {
			return true
		}
	}
	return false
}

// VowelSwapping occurs when the vowels in the target token are swapped with
// each other, leading to a slightly altered version of the original word.
// This type of error typically involves exchanging one vowel for another,
// which can still make the altered token look similar to the original,
// but with a subtle change. For example, the word "example" could become
//
//	"ixample", "exomple", or "exaple", where vowels like "a", "e", and "o"
//
// are swapped, causing the token to differ from its correct form.
// One vowel moves at a time. strings.Replace with -1 rewrote every occurrence
// at once, so a name with a repeated vowel never produced the typo anyone
// actually makes: "google" gave "gaagle" but never "gaogle" or "goagle", and
// "example" gave "ixampli" rather than the "ixample" this comment promises.
// Every case in the test suite happened to use each vowel at most once, where
// the two are the same, so nothing caught it.
//
// Positions are found by substring search rather than by rune index, because a
// language's vowel set is not always one rune per entry.
func VowelSwapping(token string, vowels ...string) []string {
	u := newUniq()
	for _, from := range vowels {
		if from == "" {
			continue
		}
		for off := 0; off < len(token); {
			i := strings.Index(token[off:], from)
			if i < 0 {
				break
			}
			at := off + i
			for _, to := range vowels {
				if to != from {
					u.add(token, token[:at]+to+token[at+len(from):])
				}
			}
			off = at + len(from)
		}
		// And the systematic substitution, which is what this function used to
		// return on its own. It is a rarer typo than a single slip but a real
		// one, and a name it generates is registrable like any other, so it is
		// kept rather than traded away for the per-position case. uniq drops
		// the duplicate when the vowel occurs once.
		for _, to := range vowels {
			if to != from {
				u.add(token, strings.ReplaceAll(token, from, to))
			}
		}
	}
	return u.tokens()
}

// HomophoneSwapping occurs when words that sound the same but have different
// meanings or spellings are substituted for one another. This type of error
// arises from words that are homophones—words that are pronounced the same but
// may differ in spelling or meaning. For example, the word "base" could be
// swapped with "bass", where "base" and "bass" are homophones, making the
// altered word sound the same when spoken, yet look different in writing.
func HomophoneSwapping(token string, homophones ...[]string) (words []string) {
	return swapWordSets(token, homophones)
}

// HomoglyphSwapping is a technique where visually similar characters, called
// homoglyphs, are swapped for one another in text. These characters look alike
// but are actually different in code, often coming from different alphabets
// or character sets. For example, an attacker might replace the letter "o" with
// the Cyrillic letter "о" (which looks nearly identical) in a URL or word. This
// can trick people into clicking a fraudulent link or misreading text.
func HomoglyphSwapping(token string, homoglyphs map[string][]string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i, r := range rs {
		for _, kchar := range similarChars(string(r), homoglyphs) {
			u.add(token, joinRunes(rs[:i], kchar, rs[i+1:]))
		}
	}
	return u.tokens()
}

// BitFlipping involves altering the binary representation of characters in a
// token by flipping one or more bits. This technique introduces subtle changes
//
//	to the characters, which can result in visually similar but distinct tokens.
//
// For example, flipping a single bit in the character "a" might produce a
//
//	different character entirely, such as "b", creating variants that are hard
//
// to detect visually but differ in encoding.
func BitFlipping(token string, graphemes ...string) (variations []string) {
	// Flip a single bit in a byte
	flipBit := func(b byte, pos uint) byte {
		mask := byte(1 << pos)
		return b ^ mask
	}

	// Flip each bit in each byte of the token
	for i := 0; i < len(token); i++ {
		for bit := 0; bit < 8; bit++ {
			flippedChar := flipBit(token[i], uint(bit))
			// Construct new variation
			variant := token[:i] + string(flippedChar) + token[i+1:]
			variations = append(variations, variant)
		}
	}
	return
}

// TokenOrderSwap involves rearranging the order of words, numbers, or components
// within a token to create alternative versions. This method often results in
// tokens that are similar to the original but with a different sequence,
// which can be used to confuse or mislead users. For example, the token
// "2024example" could be altered to "example2024", or "shop-online" could
//
//	become "online-shop", where the elements are swapped in position.
func TokenOrderSwap(token string) (variations []string) {
	parts, seps := SplitTokens(token)
	if len(parts) < 2 {
		return nil
	}

	u := newUniq()
	emit := func(order []int) {
		var b []rune
		for i, idx := range order {
			if i > 0 {
				// Separators stay in their original positions. Moving them
				// with the tokens would rewrite "shop-online24" as
				// "24-onlineshop", which is a different edit; this algorithm
				// reorders the words and nothing else.
				b = append(b, []rune(seps[i-1])...)
			}
			b = append(b, []rune(parts[idx])...)
		}
		u.add(token, string(b))
	}

	// Every permutation while there are few enough to enumerate, and only the
	// cheap ones beyond that: 5 tokens is 120 orderings and 8 is 40,320, which
	// is a combinatorial blowup for names nobody types anyway.
	if len(parts) <= 4 {
		permute(len(parts), func(order []int) { emit(order) })
		return u.tokens()
	}

	rev := make([]int, len(parts))
	for i := range rev {
		rev[i] = len(parts) - 1 - i
	}
	emit(rev)
	for i := 0; i+1 < len(parts); i++ {
		order := make([]int, len(parts))
		for j := range order {
			order[j] = j
		}
		order[i], order[i+1] = order[i+1], order[i]
		emit(order)
	}
	return u.tokens()
}

// permute calls fn with every ordering of n indices, in a deterministic order.
func permute(n int, fn func([]int)) {
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	var rec func(int)
	rec = func(k int) {
		if k == n {
			fn(order)
			return
		}
		for i := k; i < n; i++ {
			order[k], order[i] = order[i], order[k]
			rec(k + 1)
			order[k], order[i] = order[i], order[k]
		}
	}
	rec(0)
}

// SplitTokens breaks a name into the words a reader sees, and the separators
// between them.
//
// Two boundaries count: an explicit separator (- _ .) and the join between
// letters and digits, so "shop-online" is two tokens and "2024example" is also
// two. len(seps) is always len(parts)-1, and rejoining parts[i] with seps[i]
// reproduces the input exactly.
func SplitTokens(token string) (parts []string, seps []string) {
	rs := runesOf(token)
	if len(rs) == 0 {
		return nil, nil
	}

	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			parts = append(parts, string(cur))
			cur = nil
		}
	}
	isSep := func(r rune) bool { return r == '-' || r == '_' || r == '.' }
	digit := func(r rune) bool { return r >= '0' && r <= '9' }

	for i, r := range rs {
		switch {
		case isSep(r):
			flush()
			// Runs of separators belong to one boundary: "a--b" is two tokens.
			if n := len(seps); n > 0 && len(parts) == n {
				seps[n-1] += string(r)
				continue
			}
			seps = append(seps, string(r))
		default:
			if i > 0 && len(cur) > 0 && digit(r) != digit(rs[i-1]) {
				flush()
				seps = append(seps, "")
			}
			cur = append(cur, r)
		}
	}
	flush()

	// A trailing or leading separator leaves more separators than gaps; that is
	// not a token boundary and reordering around it would move the separator.
	if len(seps) != len(parts)-1 {
		return nil, nil
	}
	return parts, seps
}

// CardinalSwap involves replacing numerical digits with their corresponding
// cardinal word forms, or vice versa. This process creates variants by
// converting numbers to words or words to numbers. For example, the token
// "file2" might be altered to "filetwo", or "chapterthree" could become "chapter3".
// numeralSwap walks digit<->word substitutions in both directions and returns
// every distinct result.
//
// The seen set is shared across the whole walk rather than rebuilt in each
// frame. Per frame it only stopped a variant repeating among its own siblings,
// so any data that maps a word back to its digit sent the walk around a
// two-cycle until the stack gave out. That is not a hypothetical shape: it is
// exactly what dataset.Lang.Numerals returns, because a numeral line is stored
// as a clique of transitions and every token on it becomes a key -- "1" ->
// "first" and "first" -> "1" both exist. `urlinsane typo -a cns` crashed the
// process with a stack overflow on the shipped English numerals, while the unit
// tests passed because their fixture is keyed by digit alone and cannot cycle.
//
// Keys are walked in sorted order because ranging a map is randomised, and the
// order variants come back in reaches admission order in the engine.
func numeralSwap(token string, data map[string]string) (variations []string) {
	keys := make([]string, 0, len(data))
	for k, v := range data {
		// A degenerate entry is skipped rather than walked. strings.Replace
		// with an empty `old` does not match nothing -- it matches between
		// every character, so "abc" becomes "1a1b1c1". Each step then produces
		// a strictly longer string the seen set has never seen, the walk never
		// converges, and the process grows until it dies. An empty numeral or
		// an empty word is not a substitution anyone can make anyway.
		if k == "" || v == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	seen := map[string]bool{token: true}
	var walk func(str string, reverse bool)
	walk = func(str string, reverse bool) {
		for _, num := range keys {
			word := data[num]
			var variant string
			if reverse {
				variant = strings.Replace(str, num, word, -1)
			} else {
				variant = strings.Replace(str, word, num, -1)
			}
			if variant == str || seen[variant] {
				continue
			}
			seen[variant] = true
			variations = append(variations, variant)
			walk(variant, reverse)
		}
	}
	walk(token, false)
	walk(token, true)
	return variations
}

func CardinalSwap(token string, numerals map[string][]string) (variations []string) {
	return numeralSwap(token, numeralMap(numerals, 0))
}

// OrdinalSwap involves substituting numerical digits with their corresponding
// ordinal word forms, or converting ordinal words back into numerical digits.
// This technique generates variations by switching between numeric and
//
//	word-based representations of ordinals. For example, the token "file2" could
//	be transformed into "filesecond", or "chapterthird" might be altered to
//
// "chapter3".
func OrdinalSwap(token string, numerals map[string][]string) (variations []string) {
	return numeralSwap(token, numeralMap(numerals, 1))
}

// DotHyphenSubstitution involves substituting dots (.) with hyphens (-) or
// vice versa within a given token, creating alternative versions that resemble
// the original. This technique generates variants by interchanging these
// commonly used separators, often resulting in tokens that look similar but
// are structurally different. For example, a token like "my-example.com"
// might become "my.example.com", or "my.example-com" could be changed
// to "my-example.com".
func DotHyphenSubstitution(token string) (variations []string) {
	// Both directions. "or vice versa" is in the doc comment above and was
	// never implemented, so a name with hyphens and no dots — every package
	// and repo name shaped like "my-lib" — produced nothing at all.
	u := newUniq()
	for _, v := range characterReplace(token, ".", "-") {
		u.add(token, v)
	}
	for _, v := range characterReplace(token, "-", ".") {
		u.add(token, v)
	}
	return u.tokens()
}

// StemSwapping involves replacing words with their corresponding root or stem forms,
// or vice versa. This process generates variations by switching between the
// base form of a word and its derived forms. For example, the token "running"
// might be altered to its root "run", or "player" could become "play".
func StemSwapping(token string, tokens []string) (variations []string) {
	// TODO: create a multilingual stemmer first
	return
}

// EmojiInsertion inserts emojis in target names. This technique exploits
// the presence of emojis in the target name.
func EmojiInsertion(token string, tokens []string) (variations []string) {

	return
}
