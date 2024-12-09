# Typosquatting


Typosquatting is where attackers exploit named entities for personal or financial gain. 

They exploit typing errors and linguistics to generate named variants that trick humans into doing something they did not intened to do. This can have a wide range of consiquinces for both individuals and businesses. 




A named entity can be anything from a domain name, product name, username, software name, sourcecode libary name, and anything represented primarily by its name. 

Its a process where a similar named entity is created by an attacker to redirect traffic to it. 

It taking advantage of common human errors we make everyday and is used by cyber criminals in the following areas:

* Supply Chain Attacks
* Phishing Campaigns
* Malware Distribution
* Ad Revenue Generation
* Competitor Sabotage
* Brand Exploitation






Attackers can take a

They take advantage of human errors by variations of named entities that  


 to linguistics and topigraphical properties of the human to machine interface redicect to construct variations in ways that exploit the visual similarity or familiarity of certain languages and alphabets.

Typosquatting is a form of cybersquatting where attackers register domain names that are similar to legitimate, well-known domain names, exploiting common typographical errors or variations made by users when typing URLs into a browser.




One of the few ways to see if your named entity is being exploited is to generate a list of named variants and check if they are beings used in a malicous way. 

We will explore various algorithms used to generate variations of named entities. The type of algorytms my vary depending on the type of named entity that is targeted. For example usernames and domain names have different naming convention and constraints. 

Offten the constrains are inforce by the underlying service being used. 

Domain name standards is regulated by the Internet Corporation for Assigned Names and Numbers (ICANN) and others. 

These standards restrict domain names to follow a specificy naming format. 

For the purpose of this blog, we will focus on domains names. 


Other named entities such as product names may not have these strict global naming restrictions. 



For the purpose of exploring typo squatting we will narrow the scope to domain names only. 

We will also use the urlinsane tool to generate domain name variants and check if they are being used.











# Algorithms
Typo algoryhms exploit typing errors and linguistic patterns to generate named variants. 

These algorithms can be classified based on the way they manipulate names or tokens. 

In Natural Language Processing (NLP) a token refer to a unit of text. 

Named entities can have one or more tokens. Sometimes a single token can be split into multiple tokens in a process called to tokenization that we will examine later. 

The following are some of the categories of typosquatting algorithms that apply to domain names. 

Typo algorithms can fall into these categories: Insertion, Ommition, Transposition, and Substitution



For demo purposes we are going to focus only on domain typosquatting and use a tool called **urlinsane** that I created 


## Substitution

Algorithms in this category create named variants by replacing one or more characters or tokens in the original name.


### Adjacent Character Substitution

This algorithm generates named variants by replaces each character in the target name with one of the adjacent characters on the keyboard layout.

For example:

```sh
urlinsane typo -k en1 -a acs a
```

-k en1 : Tells the tool to use the 'en1' keyboard layout which is the for the 'QWERTY' English langauge keyboard layout. 
-a acs : Tells the tool to only use the Adjacent Character Substitution (ACS) algorithm to generate name variants
     a : Is the named entity we are tageting

Outputs:

```sh
1  ACS  q
1  ACS  w
1  ACS  s  
1  ACS  z 
```

The numbers at the begining of the line are the levinshtine distance between the named entity and the named variant.


### Grapheme Substitution

This algorithm generates named variants by replaces each character in the target name with a character in the alphabet of the given langauge. 

For example:

```sh
urlinsane typo -l en -a gs a
```
-l en : Tells the tool to use the the english language. 
-a gs : Tells the tool to only use the Grapheme Substitution (GS) algorithm to generate name variants
     a : Is the named entity we are tageting

Output (truncated):

```sh
1  GS  a
1  GS  b
1  GS  c  
...
1  GS  z 
```


### Symbols Substitution

Involves substituting dots (.) with hyphens (-) or vice versa within a given token, creating alternative versions that resemble the original. This technique generates variants by interchanging these commonly used separators, often resulting in tokens that look similar but are structurally different. For example, a token like "my-example.com" might become "my.example.com", or "my.example-com" could be changed to "my-example.com".


### Homophone Substitution 
occurs when words that sound the same but have different meanings or spellings are substituted for one another. This type of error arises from words that are homophones—words that are pronounced the same but may differ in spelling or meaning. For example, the word "base" could be swapped with "bass", where "base" and "bass" are homophones, making the altered word sound the same when spoken, yet look different in writing.


### Homoglyph Substitution 
is a technique where visually similar characters, called homoglyphs, are swapped for one another in text. These characters look alike but are actually different in code, often coming from different alphabets or character sets. For example, an attacker might replace the letter "o" with the Cyrillic letter "о" (which looks nearly identical) in a URL or word. This can trick people into clicking a fraudulent link or misreading text.

### Unicode Homoglyph Substitution



### Cardinal Substitution
involves replacing numerical digits with their corresponding cardinal word forms, or vice versa. This process creates variants by converting numbers to words or words to numbers. For example, the token "file2" might be altered to "filetwo", or "chapterthree" could become "chapter3".



### Stem Substitution 
involves replacing words with their corresponding root or stem forms, or vice versa. This process generates variations by switching between the base form of a word and its derived forms. For example, the token "running" might be altered to its root "run", or "player" could become "play".


### Ordinal Substitution 
involves substituting numerical digits with their corresponding ordinal word forms, or converting ordinal words back into numerical digits. This technique generates variations by switching between numeric and
 
 	word-based representations of ordinals. For example, the token "file2" could
 	be transformed into "filesecond", or "chapterthird" might be altered to
  "chapter3".

### Singular Plural Substitution
Thse typos are where a word is altered by switching between its singular and plural forms. This subtle change can create a word that looks similar to the original, but with a small variation that is easy to overlook.
For example:


### Misspelling Substitution
These typos are created by frequent spelling errors or missteps that occur in the target language. These errors often involve slight changes to the spelling of a word, making them appear similar to the original but incorrect. For instance, the word "youtube" could be mistyped as "youtub", and "abseil" could become "absail", where common mistakes in spelling lead to slightly altered but recognizable versions of the original.

### Multilingual Token Substitution
    For international brands, typosquatting can incorporate hybrid language variations that mix English with local language elements. For instance, combining “amazon” with country-specific words like “amazonkaufen.com” (using the German word for “buy”) can mislead German-speaking users.

### TLD Substitution




## Insertion

Insertion typo algorithms create named variants by adding extra characters or tokens to target name.


### Prefix Insertion
Adds characters to the start (e.g., `example.com` → `wwwexample.com`).

### Suffix Insertion
Adds characters to the end (e.g., `example.com` → `example-com.net`).


### Keywords Insertion

### Adjacent Character Insertion 

typos occur when characters adjacent of each letter are inserted. For example, googhle inserts "h" next to it's adjacent character "g" on an English QWERTY keyboard layout.

### Character Insertion

Also known as alphabet insertion, occurs when one or more unintended letters are added to a valid token, leading to a modified or misspelled version of the original token. These extra characters are typically inserted either at the beginning or within the token, causing it to deviate from its intended form. This type of error is often the result of a slip of the finger or an accidental keystroke. For example, the token "example" might be mistakenly typed as "aexample", "eaxample", "exaample", "examaple",
 
 	or "eaxampale", where additional letter like "a" are inserted throughout
  the token, distorting its original structure.


### Emoji Insertion 
inserts emojis in target names. This technique exploits the presence of emojis in the target name.



### Hyphen Insertion

typos happen when hyphens are mistakenly placed between characters in a token, often occurring in various positions around the letters. This type of error can lead to unnecessary fragmentation of the word, with hyphens inserted at different points throughout the token. For example:


### Dot Insertion

typos take place when periods (.) are mistakenly added at various points within a token, leading to an incorrect placement of dots in the target token. This type of error typically happens due to inadvertent key presses or misplacement while typing. For instance, the word "example" may be mistakenly written as "e.xample", "ex.ample", "exa.mple", "exam.ple", or "examp.le", where the dot appears at different locations within the token, disrupting the original structure.

### Double Character Insertion 
These are typos that occur when a letter is unintentionally repeated within a token, leading to a misspelled version. This type of error typically happens when a key is pressed twice or a letter is accidentally duplicated. 
For example: 

## Adjacent Double Character Insertion
These typos occur when consecutive, identical letters in a token are replaced with adjacent keys on the keyboard, resulting in a slight alteration of the original word. This type of error often happens due to accidental key presses of nearby characters. 
For example:



## Omission
Omission typo algorithms create named variants by removing characters or tokens from the target name.


### Dot Omission 
Typos happen when periods (.) that should be present in the target token are unintentionally omitted or left out. This type of error typically occurs when the user fails to input the expected dots, often resulting in a word or sequence that appears as a single string without proper separation. For example, the sequence "one.two.three" might be mistakenly written as "one.twothree", "onetwo.three", or even "onetwothree", where the dots are missing between certain parts of the token, causing it to lose the intended structure or meaning.


### Character Omission 
occurs when one character is unintentionally omitted from the token, leading to an incomplete version of the original word. This type of typo can happen when a key is accidentally skipped or overlooked while typing. For example, the word "google" might be mistakenly written as "gogle", "gogle", "googe", "googl", "goole", or "oogle", where a single character is missing from different positions in the word, causing it to deviate from the correct spelling.

### Double Character Omission

ie. ggle

### Hyphen Omission 
 typos occur when hyphens are unintentionally left out of a token, resulting in a version of the token that misses the expected hyphenation. For example, the token "one-for-all" might be mistakenly written as "onefor-all", "one-forall", or even "oneforall", where the hyphens are omitted.



## Transposition

Transposition typo algorithms create named variants by swapping characters or tokens in the target name.

### Character Swapping 

Refers to a type of typo where two adjacent characters in the original token are exchanged or swapped. This often occurs when characters are unintentionally reversed in order, resulting in a misspelling.For example, the word "example" could become "examlpe" by swapping the position of the letters "l" and "p".

### Vowel Swapping 
occurs when the vowels in the target token are swapped with each other, leading to a slightly altered version of the original word. This type of error typically involves exchanging one vowel for another, which can still make the altered token look similar to the original, but with a subtle change. For example:

### Token Swapping 
involves rearranging the order of words, numbers, or components within a token to create alternative versions. This method often results in tokens that are similar to the original but with a different sequence, which can be used to confuse or mislead users. For example, the token "2024example" could be altered to "example2024", or "shop-online" could
 
 	become "online-shop", where the elements are swapped in position.





## System

There are other ways named variants are created at no fault of the human operator. These are created by system or software errors.

### Bit Flipping

BitFlipping involves altering the binary representation of characters in a token by flipping one or more bits. This technique introduces subtle changes to the characters, which can result in visually similar but distinct tokens.
For example:





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
    Attackers take advantage of languages with characters that look similar to those in the target language (usually English). For instance, many Cyrillic letters closely resemble Latin ones, making it possible to create domains like “fаcebook.com” (using Cyrillic "а" instead of Latin "a") that look almost identical to the legitimate "facebook.com."    Other language scripts, such as Greek, also have similar characters, enabling a wide range of homograph attacks, which exploit look-alike characters from different languages.

2. Misspellings and Common Typing Errors
    Typo domains rely on common typing errors, especially for language-specific keyboards. For instance, a Spanish-speaking audience might accidentally type “goolge” instead of “google” due to familiarity with specific keystroke patterns.    Attackers might also replace accented letters or letters with diacritics in languages like Spanish, French, or German. For example, using "télégramme.com" for “telegram.com” could mislead French-speaking users.

3. Phonetic and Linguistic Variations
    Attackers exploit phonetic similarities, where the target name is replaced by a similar-sounding word or phrase in the target language. For instance, “secure” might be replaced by “sekur” or “bank” by “banq” for a French-speaking audience.    Additionally, attackers create typosquatted domains that reflect common dialectal or regional language spellings, targeting specific communities. For instance, “colour.com” might be typosquatted as “colur.com” to confuse UK users.

4. Multilingual Blending
    For international brands, typosquatting can incorporate hybrid language variations that mix English with local language elements. For instance, combining “amazon” with country-specific words like “amazonkaufen.com” (using the German word for “buy”) can mislead German-speaking users.

5. Foreign Language Cognates and Loanwords
    Attackers use cognates (words that look similar and have the same meaning across languages) or loanwords to appeal to international audiences. A word like “hotel” appears in many languages and can be combined with a typosquatting element like “h0tel.com” to fool non-native English speakers, as they may overlook minor changes in familiar words.

6. Transliterations and Alternate Alphabets
    Some typosquatters use transliterations, converting words from one alphabet to another in ways that resemble the original brand or name. For example, using Punycode (a way of encoding Unicode within the ASCII character set) allows for domain names like “xn--pple-43d.com,” which appears as “аpple.com” in the browser (with Cyrillic "а").    Transliteration is also common in languages like Arabic, Chinese, and Hindi, where brand names are spelled out using Latin characters or phonetically similar sounds to trick users.

7. Local Language Targeting
    Typosquatting often uses localized or region-specific spellings and slang to make a fake domain appear legitimate to a particular audience. For instance, Spanish speakers might see “amig0s.com” instead of “amigos.com,” where “0” is used in place of “o” to fool users who are accustomed to similar regional variations.

By carefully crafting typosquatted names that resonate linguistically and culturally with the target audience, attackers enhance the believability of their fake domains or usernames, increasing the likelihood of successful phishing or other deceptive attacks.








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










