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


Here's a structured summary of the plugin types and their roles in the application:
| **Type**    | Count  | Decscription |
|-------------|--------|--------------|
| Languages   |    9   | Language plugins enable support for various language models, expanding the application's linguistic capability. |
| Keyboards   |    19  | Keyboard layout plugins allow us to target multiple languages and regions |
| Algorithms  |    24  | Used to generate typo variants for domains, arbitrary names, and software packages|
| Information |    8   | Used for gather data on domains, software libraries, and named entities|
| Outputs     |    6   | Formats data for display and or save outputs to files, improving usability and reporting |
| Database    |    1   | It caches and saves scan results and boosting performance and enabling efficient data retrieval.|



# Languages

In typosquatting, language plays a significant role in manipulating legitimate terms and names to create deceptive variations that appear familiar to the target audience. Attackers use linguistic techniques to construct these variations in ways that exploit the visual similarity or familiarity of certain languages and alphabets.

| Languages | Keyboards | Homoglyphs | Homophones | Antonyms | Misspellings | Cardinal | Ordinal | Vowels | Graphemes  | 
|-----------|-----------|-----------|------------|-----------|--------------|--------|-----------|---------|-----------|
| Arabic    |    4      |           |            |           |              |   |   |   |    | 
| Armenian  |    3      |           |            |           |              |   |   |   |    | 
| English   |    4      |           |            |           |              |   |   |   |    | 
| Finnish   |    1      |           |            |           |              |   |   |   |    | 
| Russian   |    3      |           |            |           |              |   |   |   |    | 
| Spanish   |    2      |           |            |           |              |   |   |   |    | 
| Hebrew    |    1      |           |            |           |              |   |   |   |    | 
| Persian   |    1      |           |            |           |              |   |   |   |    | 

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



## Algorithms

| ID | Name   | Description  | 
|--|------------|-------------|
| |   |         |    
| |    |           |   
| |    |         |   


| ID | Name            | Description ||
|---|---------------------------------|--|
| | Missing Dot                     ||
| | Missing Dashes                  | |
| | Strip Dashes                    | |
| | Character Omission              | |
| | Character Swap                  | |
| | Adjacent Character Substitution | |
| | Adjacent Character Insertion    ||
| | Homoglyphs                      ||
| | Singular Pluralise              | |
| | Character Repeat                | |
| | Double Character Replacement    | |
| | Common Misspellings             | |
| | Homophones                      ||
| | Vowel Swapping                  | |
| | Bitsquatting                    | |
| | Wrong Top Level Domain          | | 
| | Wrong Second Level Domain       | | 
| | Wrong Third Level Domain        | |
| | Ordinal Number Swap             | |
| | Cardinal Number Swap            ||
| | Hyphenation                     || 
| | Multithreaded Algorithms        ||   
| | Subdomain insertion             | |
| | Period Insertion                | | 
| | Combosquatting (Keywords)       | |



## Information

| ID | Name   | Description  | 
|----|------------|-------------|
|    |   |         |    
|    |    |           |   
|    |    |         |   

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






