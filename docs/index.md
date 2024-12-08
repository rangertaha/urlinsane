# URLInsane

Generates domain typos and variations used to detect and perform typo squatting, 
URL hijacking, phishing, and corporate espionage.  Inspired by URLCrazy I wanted 
to create a better version that supported multiple languages and linguistic typos.
I also wanted it to be a binary with fast execution time.



Table of contents
=================

<!--ts-->
   * [Table of contents](#table-of-contents)
   * [Introduction](#introduction)
   * [Installation](#installation)
   * [Usage](#usage)
   * [Features](#features)
   * [Languages](#languages)
      * [English](#english)
      * [Spanish](#spanish)
      * [Russian](#russian)
      * [Finish](#finish)
      * [Arabic](#arabic)
      * [Persian](#persian)
      * [Hebrew](#hebrew)
   * [Algorithms](#algorithms)
   * [Extra Functions](#extra-functions)
      * [TODO](#todo)
   * [Authors](#authors)
   * [License](#license)
<!--te-->





## Introduction
Generate and test domain typos and variations to detect and perform typo squatting, URL hijacking, phishing, and corporate espionage.

The engine is designed to execute concurrent typo algorithms and then additional 
concurrent functions for each domain variation. The additional functions can 
check DNS records and much more. It's also designed for extensibility, allowing 
developers to add functionality and support for additional languages. See 
[URLInsane](https://rangertaha.github.io/urlinsane/) for more details.


# URLInsane


[![Go Report Card](https://goreportcard.com/badge/github.com/rangertaha/urlinsane?style=flat-square)](https://goreportcard.com/report/github.com/rangertaha/urlinsane) [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/rangertaha/urlinsane) [![PkgGoDev](https://pkg.go.dev/badge/github.com/rangertaha/urlinsane)](https://pkg.go.dev/github.com/github.com/rangertaha/urlinsane) [![Release](https://img.shields.io/github/release/rangertaha/urlinsane.svg?style=flat-square)](https://github.com/rangertaha/urlinsane/releases/latest) [![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)



URLInsane is a robust command-line tool designed to detect typosquatting across domains, arbitrary names, usernames, and software packages. By leveraging advanced algorithms, information-gathering techniques, and data analysis, it identifies potentially harmful variations of targeted entities that cybercriminals might exploit. Essential for defending against typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, and corporate espionage, URLInsane also enhances threat intelligence capabilities.

Featuring a plugin-based multilingual permutation engine, URLInsane supports various keyboard layouts for multiple languages. Its extensible plugin system allows easy addition of new capabilities. Currently, it includes plugins for 24 algorithms, 8 information gathering methods, 19 keyboard layouts, 9 languages, and 4 output formats. Originally developed in Python, URLInsane was built in Go to enhance speed, concurrency, and portability.



Urlinsane is a powerful command-line tool designed for the detection of typo-squatting across domains, email addresses, usernames, and software packages. Utilizing  Urlinsane scans for potential typosquatted variants by applying advanced typo squatting algorithms, information gathering, and data analysis.  It identifies potentially harmful variations of a victim's domain name, email address, software packages, and that cybercriminals might exploit. 



Urlinsane is used to aid in the detection of typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage, and threat intelligence.

It's a plugin-based multilingual permutation engine that supports keyboard layouts for each language. The plugin system allows us to easily extend its capabilities. Currently, it supports plugins for algorithms (24), information gathering (8), keyboard layouts (19), languages (9), and output (4) formats. Originally, I wrote similar tools in Python but was not happy with the performance. This tool was built in Golang to take advantage of its speed, concurrency, and portability (see [`Speed`](#Speed).

# Features
* Fast execution time
* Cuncurrency in generating typos and collection informaiton
* Distribute queries to multiple DNS Servers
* Single downloadable binary with no system depedencies
* Plugin architecture for extensability
* Multi-lingual language modeling



# Plugins

Plugins play a crucial role in extending the functionality, flexibility, and customization of Urlinsane and allow it to evolve alongside changing needs and technological advancements.


Here's a structured summary of the plugin types and their roles in Urlinsane:
| **Type**    | Count  | Decscription |
|-------------|--------|--------------|
| Languages   |    9   | Language plugins enable support for various language models, expanding the application's linguistic capability. |
| Keyboards   |    19  | Keyboard layout plugins allow us to target multiple languages and regions |
| Algorithms  |    24  | Used to generate typo variants for domains, arbitrary names, and software packages|
| Information |    8   | Used for gather data on domains, software libraries, and named entities|
| Outputs     |    6   | Formats data for display and or save outputs to files, improving usability and reporting |
| Database    |    1   | It caches and saves scan results and boosting performance and enabling efficient data retrieval.|



# Languages

In typosquatting, language plays a significant role in manipulating legitimate tokens or names to create deceptive variations that appear familiar to the target audience. Attackers exploit linguistic, topigraphical properties of human machine interface to construct these variations in ways that exploit the visual similarity or familiarity of certain languages and alphabets.

## Keyboards
## Homoglyphs
## Homophones
## Antonyms
## Misspellings 
## Cardinal 
## Ordinal 
## Vowels
## Graphemes  





1. Homoglyph Replacement Across Alphabets

    Attackers take advantage of languages with characters that look similar to those in the target language (usually English). For instance, many Cyrillic letters closely resemble Latin ones, making it possible to create domains like “fаcebook.com” (using Cyrillic "а" instead of Latin "a") that look almost identical to the legitimate "facebook.com."
    Other language scripts, such as Greek, also have similar characters, enabling a wide range of homograph attacks, which exploit look-alike characters from different languages.

2. Misspellings and Common Typing Errors

    Typo domains rely on common typing errors, especially for language-specific keyboards. For instance, a Spanish-speaking audience might accidentally type “goolge” instead of “google” due to familiarity with specific keystroke patterns.
    Attackers might also replace accented letters or letters with diacritics in languages like Spanish, French, or German. For example, using "télégramme.com" for “telegram.com” could mislead French-speaking users.

3. Phonetic and Linguistic Variations

    Attackers exploit phonetic similarities, where the target name is replaced by a similar-sounding word or phrase in the target language. For instance, “secure” might be replaced by “sekur” or “bank” by “banq” for a French-speaking audience.
    Additionally, attackers create typosquatted domains that reflect common dialectal or regional language spellings, targeting specific communities. For instance, “colour.com” might be typosquatted as “colur.com” to confuse UK users.

4. Multilingual Blending

    For international brands, typosquatting can incorporate hybrid language variations that mix English with local language elements. For instance, combining “amazon” with country-specific words like “amazonkaufen.com” (using the German word for “buy”) can mislead German-speaking users.

5. Foreign Language Cognates and Loanwords

    Attackers use cognates (words that look similar and have the same meaning across languages) or loanwords to appeal to international audiences. A word like “hotel” appears in many languages and can be combined with a typosquatting element like “h0tel.com” to fool non-native English speakers, as they may overlook minor changes in familiar words.

6. Transliterations and Alternate Alphabets

    Some typosquatters use transliterations, converting words from one alphabet to another in ways that resemble the original brand or name. For example, using Punycode (a way of encoding Unicode within the ASCII character set) allows for domain names like “xn--pple-43d.com,” which appears as “аpple.com” in the browser (with Cyrillic "а").
    Transliteration is also common in languages like Arabic, Chinese, and Hindi, where brand names are spelled out using Latin characters or phonetically similar sounds to trick users.

7. Local Language Targeting

    Typosquatting often uses localized or region-specific spellings and slang to make a fake domain appear legitimate to a particular audience. For instance, Spanish speakers might see “amig0s.com” instead of “amigos.com,” where “0” is used in place of “o” to fool users who are accustomed to similar regional variations.

By carefully crafting typosquatted names that resonate linguistically and culturally with the target audience, attackers enhance the believability of their fake domains or usernames, increasing the likelihood of successful phishing or other deceptive attacks.




# Algorithms


Typosquatting algorithms can be categorized depending on how they manipulate target names



## Substitution

Algorithms in this category create typos by replacing one or more characters in the original name



- **Neighboring Key Substitution**: Swaps characters with those near them on a keyboard (e.g., `example.com` → `exzmple.com`).
- **Visual Substitution (Homoglyphs)**: Replaces characters with similar-looking alternatives (e.g., `google.com` → `g00gle.com`).
- **TLD Substitution**: Uses a different valid TLD (e.g., `.com` → `.net`, `example.com` → `example.net`).
- **Unicode Homoglyphs**: Substitutes characters with Unicode equivalents (e.g., `google.com` → `gοοgle.com` where `ο` is Greek).
- **Brand Misspelling**: Slightly alters the brand name (e.g., `facebook.com` → `facebok.com`).


### Adjacent Character Substitution

These typos happen when a character in the original token is mistakenly replaced by a neighboring character from the same keyboard layout. This type of error often occurs due to hitting an adjacent key by accident.
For example:


// GraphemeReplacement, also known as alphabet replacement, occurs when characters
// from the original token are replaced by other letters from the alphabet,
// resulting in a modified version of the token. This type of error typically leads
// to small changes in the original token, where one or more letters are swapped
// for different characters. For example, the token "example" could be mistakenly
// written as "axample", "bxample", "cxample", "dxample", or "eaample", where
// letters like "a", "b", "c", "d", or "e" are substituted, altering the
// word slightly but keeping its general structure.


// DotHyphenSubstitution involves substituting dots (.) with hyphens (-) or
// vice versa within a given token, creating alternative versions that resemble
// the original. This technique generates variants by interchanging these
// commonly used separators, often resulting in tokens that look similar but
// are structurally different. For example, a token like "my-example.com"
// might become "my.example.com", or "my.example-com" could be changed
// to "my-example.com".



## Insertion
These types of typo algorithms add extra characters or tokens to target name.


- **Single Character Insertion**: Adds a single character anywhere in the domain (e.g., `example.com` → `examp1e.com`).
- **Double Insertion**: Adds two or more characters (e.g., `example.com` → `exampllee.com`).
- **Hyphenation**: Adds or removes hyphens (e.g., `example.com` → `ex-ample.com`).
- **Prefix Insertion**: Adds characters to the start (e.g., `example.com` → `wwwexample.com`).
- **Suffix Addition**: Adds characters to the end (e.g., `example.com` → `example-com.net`).
- **Subdomain Addition**: Adds subdomains to mimic legitimate structures (e.g., `example.com` → `login.example.com`).
- **Keyword Insertion**: Adds popular or generic terms (e.g., `example.com` → `example-news.com`).
- **Phrase Insertion**: Merges multiple brand-like terms (e.g., `example.com` → `bestexample.com`).
- **Brand Affiliation**: Adds keywords associated with the brand (e.g., `facebook.com` → `facebooklogin.com`).


// AdjacentCharacterInsertion typos occur when characters adjacent of each
// letter are inserted. For example, googhle inserts "h" next to it's
// adjacent character "g" on an English QWERTY keyboard layout.

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


// EmojiInsertion inserts emojis in target names. This technique exploits
// the presence of emojis in the target name.



// HyphenInsertion typos happen when hyphens are mistakenly placed between
// characters in a token, often occurring in various positions around the
// letters. This type of error can lead to unnecessary fragmentation of the
// word, with hyphens inserted at different points throughout the token.
// For example, the word "example" might be incorrectly written as "-example",
//
//	"e-xample", "ex-ample", "exa-mple", "exam-ple", "examp-le", or even
//
// "example-", with hyphens appearing before, between, or after the letters.

// DotInsertion typos take place when periods (.) are mistakenly added at
// various points within a token, leading to an incorrect placement of dots in
// the target token. This type of error typically happens due to inadvertent
// key presses or misplacement while typing. For instance, the word "example"
// may be mistakenly written as "e.xample", "ex.ample", "exa.mple", "exam.ple",
// or "examp.le", where the dot appears at different locations
// within the token, disrupting the original structure.




## Omission
These types of typo algorithms involves removing characters or tokens from the target name.


- **Single Character Omission**: Deletes one character (e.g., `example.com` → `exaple.com`).
- **Double Omission**: Deletes two or more characters (e.g., `example.com` → `exmple.com`).
- **Hyphenation**: Adds or removes hyphens (e.g., `example.com` → `ex-ample.com`).



// DotOmission typos happen when periods (.) that should be present in the target
// token are unintentionally omitted or left out. This type of error typically
// occurs when the user fails to input the expected dots, often resulting in a
// word or sequence that appears as a single string without proper separation.
// For example, the sequence "one.two.three" might be mistakenly written
// as "one.twothree", "onetwo.three", or even "onetwothree", where the dots
// are missing between certain parts of the token, causing it to lose the
// intended structure or meaning.



// CharacterOmission occurs when one character is unintentionally omitted from
// the token, leading to an incomplete version of the original word. This type
// of typo can happen when a key is accidentally skipped or overlooked while
// typing. For example, the word "google" might be mistakenly written as "gogle",
// "gogle", "googe", "googl", "goole", or "oogle", where a single character is
// missing from different positions in the word, causing it to deviate from
// the correct spelling.



// HyphenOmission typos occur when hyphens are unintentionally left out of a
// token, resulting in a version of the token that misses the expected hyphenation.
// For example, the token "one-for-all" might be mistakenly written as "onefor-all",
// "one-forall", or even "oneforall", where the hyphens are omitted.
func HyphenOmission(token string) (tokens []string) {
	return characterDeletion(token, "-")
}




## Transposition
This types of algorithms swap characters or tokens in the target name.

### Character Swapping 

Refers to a type of typo where two adjacent characters in the original token are exchanged or swapped. This often occurs when characters are unintentionally reversed in order, resulting in a misspelling.For example, the word "example" could become "examlpe" by swapping the position of the letters "l" and "p".

/ VowelSwapping occurs when the vowels in the target token are swapped with
// each other, leading to a slightly altered version of the original word.
// This type of error typically involves exchanging one vowel for another,
// which can still make the altered token look similar to the original,
// but with a subtle change. For example, the word "example" could become
//
//	"ixample", "exomple", or "exaple", where vowels like "a", "e", and "o"
//
// are swapped, causing the token to differ from its correct form.

// HomophoneSwapping occurs when words that sound the same but have different
// meanings or spellings are substituted for one another. This type of error
// arises from words that are homophones—words that are pronounced the same but
// may differ in spelling or meaning. For example, the word "base" could be
// swapped with "bass", where "base" and "bass" are homophones, making the
// altered word sound the same when spoken, yet look different in writing.


// HomoglyphSwapping is a technique where visually similar characters, called
// homoglyphs, are swapped for one another in text. These characters look alike
// but are actually different in code, often coming from different alphabets
// or character sets. For example, an attacker might replace the letter "o" with
// the Cyrillic letter "о" (which looks nearly identical) in a URL or word. This
// can trick people into clicking a fraudulent link or misreading text.



// TokenOrderSwap involves rearranging the order of words, numbers, or components
// within a token to create alternative versions. This method often results in
// tokens that are similar to the original but with a different sequence,
// which can be used to confuse or mislead users. For example, the token
// "2024example" could be altered to "example2024", or "shop-online" could
//
//	become "online-shop", where the elements are swapped in position.


// CardinalSwap involves replacing numerical digits with their corresponding
// cardinal word forms, or vice versa. This process creates variants by
// converting numbers to words or words to numbers. For example, the token
// "file2" might be altered to "filetwo", or "chapterthree" could become "chapter3".



// StemSwapping involves replacing words with their corresponding root or stem forms,
// or vice versa. This process generates variations by switching between the
// base form of a word and its derived forms. For example, the token "running"
// might be altered to its root "run", or "player" could become "play".





// OrdinalSwap involves substituting numerical digits with their corresponding
// ordinal word forms, or converting ordinal words back into numerical digits.
// This technique generates variations by switching between numeric and
//
//	word-based representations of ordinals. For example, the token "file2" could
//	be transformed into "filesecond", or "chapterthird" might be altered to
//
// "chapter3".



## Repetition
These types of algorithms repeat one or more characters or tokens in the given name.

### Character Repetition 
These are typos that occur when a letter is unintentionally repeated within a token, leading to a misspelled version. This type of error typically happens when a key is pressed twice or a letter is accidentally duplicated. 
For example: 

## Adjacent Character Repetition
These typos occur when consecutive, identical letters in a token are replaced with adjacent keys on the keyboard, resulting in a slight alteration of the original word. This type of error often happens due
// to accidental key presses of nearby characters. 
For example:




## Linguistics

Linuistic typo algorithms exploit language tokens in the given name.


### Singular Pluralise 
Thse typos are where a word is altered by switching between its singular and plural forms. This subtle change can create a word that looks similar to the original, but with a small variation that is easy to overlook.
For example:


## Common Misspellings 
These typos are created by frequent spelling errors or missteps that occur in the target language. These errors often involve slight
// changes to the spelling of a word, making them appear similar to the original
// but incorrect. For instance, the word "youtube" could be mistyped as
// "youtub", and "abseil" could become "absail", where common mistakes in
// spelling lead to slightly altered but recognizable versions of the original.




## Digital

Digital typo algorithms exploit computer hardware to generate typos in target name.


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










# Information

Generating typos is not enought to identify posible typosquatting. We need to collect additional information on the variant that allows us to compare. 



* **IDNA**  Show international domain name (Default)
* **MX**    Checking for DNS's MX records
* **TXT**   Checking for DNS's TXT records
* **IP**    Checking for IP address
* **NS**    Checks DNS NS records
* **CNAME** Checks DNS CNAME records
* **SIM**   Show domain similarity % using fuzzy hashing with ssdeep
* **LIVE**  Show domains with IP addresses only
* **301**   Show domains redirects
* **GEO**   Show country location of IP address

## Information Gathering

|  Name           | Description  |
|-------------------------------------||
| DNS A Records                       | Retrieving IPv4 and IPv6 IP host addresses |
| DNS MX Records                      |Retrieving Mail Exchange (MX) records|
| DNS TXT Records                     |Retrieving TXT records storing arbitrary data associated with a domain |
| DNS AAAA Records                    ||
| DNS CName Records                   ||
| DNS NS Records                      |Checks DNS NS records |
| Geographic Info                     | Show country location of IP address|
| Domain Similarity                   | Show domain similarity % using fuzzy hashing with ssdeep|
| Domain Redirects                    |Show domains redirects |
| IDNA Format                         |Show international domain name (Default) |
| HTTP/SMTP Banner                    | |


DNS TXT records (Text Records) are a type of DNS record used to store arbitrary text data associated with a domain. Originally intended for descriptive text, they’re now widely used for various purposes, including domain verification, email authentication, and configuration data. TXT records allow domain owners to associate key-value data with their domain, which can be retrieved by external systems for verification and configuration purposes.

DNS MX records (Mail Exchange records) are a type of DNS record used to route email for a domain to designated mail servers. They help direct emails sent to a domain (e.g., user@example.com) to the correct mail servers that handle receiving and processing the email.



## Outputs

| Name  | Description | 
|-------|-------------|
| TABLE |         |  
| TEXT  |         |  
| CSV   |         |    
| TSV   |         |   
| MD    |         |   


## Database

| Name   | Description  | 
|--------|-------------|
| Badger |         |    












### Speed

| **Tool**   | google.com  | facebook.com  | youtube.com   | amazon.com | amazon4you.com |
|------------|-------------|---------------|---------------|------------|----------------|
| URLInsane  |             |               |               |            |                |
| URLCrazy   |             |               |               |            |                |
| DNSTwist   |             |               |               |            |                |
| DomainFuzz |             |               |               |            |                |

## Authors

* [Rangertaha (rangertaha@gmail.com)](https://github.com/rangertaha)

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details






