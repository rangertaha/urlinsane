# Typo Algorithms


## (cs) Character Swapping

Character swapping in typosquatting is a technique where two consecutive 
characters in a domain name are swapped to create a typo that looks similar to 
the original domain. This subtle change can trick users into visiting the wrong 
website, especially if they are scanning quickly and not looking closely at 
each letter.

    Example: www.examlpe.com instead of www.example.com


Character swapping takes advantage of the fact that minor typographical errors 
often go unnoticed. By switching two adjacent characters, attackers create 
domains that look nearly identical to the intended site but are different 
enough to lead users to a fake or malicious site. Example:

    examlpe.com (swapping "l" and "p")
    exapmle.com (swapping "p" and "m")


```bash
urlinsane typo -a cs -d example.com
```

```bash
urlinsane typo -a cs  -n username
```

```bash
urlinsane typo -a cs  -e username@example.com
```
- **-a** or **--algorithms** Allows you to specify the type of typosquatting algorithms to apply for generating variants
- **-d** or **--domain** Allows us to specify the target type(domain) followed by the name.


algo "github.com/rangertaha/urlinsane/pkg/typo"

## (acs) Adjacent Character Substitution

Adjacent Character Substitution typos replace characters in the original domain name with neighboring characters on a specific keyboard layout, such as QWERTY.

    Example: www.ezample.com (substitutes "z" for "x" due to proximity on a QWERTY keyboard)





## (aci) Adjacent Character Insertion

Adjacent Character Insertion typos involve adding characters adjacent to each letter on the keyboard.

    Example: www.googhle.com (inserts "h" next to "g" on a QWERTY keyboard)

## (hi) Hyphen Insertion

Hyphen Insertion typos add hyphens in a domain name at different points to create variations.

    Examples: -example, e-xample, ex-ample, exa-mple, exam-ple, examp-le, example-

## (ho) Hyphen Omission

## Hyphen Omission typos involve removing hyphens from domain names that would normally contain them.

    Example: my-example.com becomes myexample.com







## (pi) Dot Insertion

Dot Insertion typos occur when dots (.) are inserted into different parts of the domain name.

    Examples: e.xample.com, ex.ample.com, exa.mple.com, exam.ple.com



## (do) Dot Omission

Dot Omission typos leave out dots that are part of the domain.

    Example: one.two.three.com might become onetwo.three.com or one.twothree.com

## (gi) Grapheme Insertion

Grapheme Insertion (or Alphabet Insertion) adds extra letters to the original domain name to create a slight variation.

    Examples: aexample.com, bexample.com, cexample.com, dexample.com

## (gr) Grapheme Replacement

Grapheme Replacement (or Alphabet Replacement) involves replacing characters in the domain with different alphabet letters.

    Examples: axample.com, bxample.com, cxample.com, dxample.com

## (cr) Character Repetition

Character Repetition typos repeat a letter in the domain name.

    Examples: eexample.com, exaample.com, exammple.com, examplee.com

## DoubleCharacterAdjacentReplacement

Double Character Adjacent Replacement typos replace consecutive identical letters with adjacent keys on the keyboard.

    Examples: gppgle.com and giigle.com (replacing "oo" with adjacent keys on a QWERTY layout)

## (co) Character Omission

Character Omission typos occur when one character is left out from the domain name.

    Examples: gogle.com, googe.com, googl.com

## (sps) Singular Pluralise Substitution

Singular-Plural Substitution is when singular forms of words are swapped for plural forms (or vice versa) in a domain.

    Examples: example.com becomes examples.com, or examples.com becomes example.com

## Character Deletion

Character Deletion is similar to Character Omission but usually involves removing multiple characters to create a shortened version.

    Example: example.com might become exampl.com or xample.com

## (cm) Common Misspellings

Common Misspelling typos involve using frequent misspellings of words or brand names.

    Examples: youtube.com becomes youtub.com, or abseil.com becomes absail.com

## (vs) Vowel Swapping

Vowel Swapping replaces vowels in the domain name with other vowels to create variations.

    Examples: example.com might become ixample.com, exemple.com, exomple.com

## (hs) Homophone Substitution

Substitute words that sound the same but have different spellings.

    Examples: base.com becomes bass.com, site.com might become sight.com

## (hr) Homoglyph Replacement

Homoglyph substitution replaces characters with visually similar ones from different alphabets or character sets.

    Example: google.com might be replaced with googIe.com (using a capital "I" instead of a lowercase "l")

## TopLevelDomain

Top-Level Domain (TLD) Replacement changes the TLD of a domain to a similar or common alternative.

    Examples: example.com might become example.net or example.co

## SecondLevelDomain

Second-Level Domain Replacement changes the second-level part of the domain name (the main part) with similar-looking or related words.

    Example: google.com might become gogle.com or goog1e.com

## ThirdLevelDomain

Third-Level Domain Replacement involves manipulating the subdomain part of the URL.

    Example: blog.example.com might become bl0g.example.com

## (bf) BitFlipping

Bit Flipping is a low-level manipulation where individual bits in the binary representation of a domain name are flipped, creating similar-looking domains.

    Example: example.com might become exampIe.com (flipping a bit to make the "l" into a capital "I")

## (cns) Cardinal Numeral Substitution

Substitute letters with those that look similar in specific fonts, such as "1" and "l" or "0" and "O".

    Example: google.com might become goog1e.com

## (ons) Ordinal Numeral Substitution

Ordinal Swapping involves rearranging letters within the domain to form typos.

    Example: example.com might become exmaple.com

These techniques help attackers create URLs that look legitimate, making it easier to deceive users and conduct phishing attacks.

## (dh) DHSubstitution   
Dot Hyphen Substitution


## (dhu) DHUSubstitution  Dot Hyphen Underscore Substitution 
These typos are created by changing a dot to a dash.



## Stem Substitution 

Replacing words with the root form  is the process of reducing inflected (or sometimes derived) words to their word stem, base or root form



# Typo Algorithms


## CharacterSwap
characterSwapFunc typos are when two consecutive characters are swapped in the original domain name.
Example: www.examlpe.com


## AdjacentCharacterSubstitution

Adjacent character substitution typos occur when characters in the original
domain name are replaced by neighboring characters on a specific keyboard
layout. For example, www.ezample.com uses "z" instead of "x," substituting it
with the adjacent character on a QWERTY keyboard.



## AdjacentCharacterInsertion

Adjacent character insertion typos occur when characters adjacent to each
letter are inserted. For example, www.googhle.com inserts "h" next to its
adjacent character "g" on a QWERTY keyboard.





## HyphenInsertion
Hyphen insertion typos occur when hyphens are inserted adjacent to each
letter in a name. For example: "-example", "e-xample", "ex-ample", "exa-mple",
"exam-ple", "examp-le", "example-"




## HyphenOmission

## DotInsertion
Dot insertion typos occur when dots(.) are inserted the target name
For example: "e.xample", "ex.ample", "exa.mple", "exam.ple", "examp.le"


## DotOmission
Dot ommission typos occur when dots(.) are left out of the target name
For one.two.three: "one.twothree", "onetwo.three", "onetwothree",



## GraphemeInsertion
Grapheme insertion also known as alphabet insertion where additional
letters are inserted into a legitimate name to create a slightly modified
version. For example: "aexample", "bexample", "cexample", "dexample", "eaxample"



## GraphemeReplacement
Grapheme replacement also known as alphabet replacement is where additional
characters from the alphabet are replaced with characters from the target name
to produce slightly modified version. For example: "axample", "bxample",
"cxample", "dxample", "eaample"


## CharacterRepetition
Character repetition typos are created by repeating a letter in the name.
For example: "eexample", "exaample", "exammple", "examplee", "examplle"


## DoubleCharacterAdjacentReplacement

DoubleCharacterAdjacentReplacement
Double character adjacent replacement typos are created by replacing identical,
consecutive letters in the name with adjacent keys on the keyboard.
For example, www.gppgle.com and www.giigle.com.
Example keyboard layout
//
//	var keyboard = []string{
//		"1234567890-",
//		"qwertyuiop ",
//		"asdfghjkl  ",
//		"zxcvbnm    ",
//	}
//

## CharacterOmission
Grapheme omission leaves out one character from the name.
For google: "gogle", "gogle", "googe", "googl", "goole", "oogle",



The technique of creating typosquatting domains by switching between singular
and plural forms of words are often referred to as Singular-Plural Substitution
or Singular-Plural Manipulation.

## SingularPluraliseSubstitution
For instance, if the original domain is 'example', a Singular-Plural
Substitution typo might be 'examples', or vice versa. This subtle variation
can make the fake domain look credible, especially when users are quickly
scanning URLs.


## CharacterDeletion




## CommonMisspellings
Created from  common misspellings in the given language.
For example, www.youtube.com becomes www.youtub.com and www.abseil.com
becomes www.absail.com



## VowelSwapping
Created from vowels of the target name
For example,



## HomophoneSwapping
homophonesFunc are created from sets of words that sound the same when spoken.
For example, www.base.com becomes www .bass.com.



## HomoglyphSwapping
Homoglyph swapping is a technique where visually similar characters, called
homoglyphs, are swapped for one another in text. These characters look alike
but are different in code, often coming from different alphabets
or character sets. For example, an attacker might replace the letter "o" with
the Cyrillic letter "о" (which looks nearly identical) in a URL or word. This
can trick people into clicking a fraudulent link or misreading text.


## TopLevelDomain




## SecondLevelDomain



## ThirdLevelDomain



## BitFlipping


## CardinalSwapping

## OrdinalSwapping
