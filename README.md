# URLInsane

[![Go Report Card](https://goreportcard.com/badge/github.com/rangertaha/urlinsane?style=flat-square)](https://goreportcard.com/report/github.com/rangertaha/urlinsane) [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/rangertaha/urlinsane) [![PkgGoDev](https://pkg.go.dev/badge/github.com/rangertaha/urlinsane)](https://pkg.go.dev/github.com/github.com/rangertaha/urlinsane) [![Release](https://img.shields.io/github/release/rangertaha/urlinsane.svg?style=flat-square)](https://github.com/rangertaha/urlinsane/releases/latest) [![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)

Urlinsane is used to aid in the detection of typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage, supply chain attacks, and threat intelligence. It's a command-line tool for detecting typosquatting domains. It scans for potential typosquatting variants by applying advanced typo squatting algorithms, information gathering, and data analysis.  It identifies potentially harmful variations of a victim's domain name that cybercriminals might exploit. 

It's inspired by [URLCrazy](https://morningstarsecurity.com/research/urlcrazy), [Dnstwist](https://github.com/elceef/dnstwist), [DomainFuzz](https://github.com/monkeym4ster/DomainFuzz) and a few other libraries and tools I was researching at the time.





## Installation

This tools is primarly ment for Linux operating systems. 

* [urlinsane-0.6.0-darwin-amd64](https://github.com/rangertaha/urlinsane/releases/download/0.6.0/urlinsane-0.6.0-darwin-amd64)
* [urlinsane-0.6.0-linux-amd64](https://github.com/rangertaha/urlinsane/releases/download/0.6.0/urlinsane-0.6.0-linux-amd64)
* [urlinsane-0.6.0-windows-amd64.exe](https://github.com/rangertaha/urlinsane/releases/download/0.6.0/urlinsane-0.6.0-windows-amd64.exe)




## Usage

```bash
urlinsane typo example.com 
```



## Plugins

Plugins play a crucial role in extending the functionality, flexibility, and customization of Urlinsane and allow it to evolve alongside changing needs and technological advancements. Here's a structured summary of the plugin types and their roles in Urlinsane:

|    Type     | Number  | Description |
|-------------|--------|--------------|
| Languages   |    9   | Language plugins provide data for it's linguistic capability. |
| Keyboards   |    19  | Keyboard plugins provide layouts for international keyboads |
| Algorithms  |    24  | They generate typo variants for each target domain |
| Information |    13  | Collects information on target domains |
| Outputs     |    6   | Formats and saves results  |


### Languages

In typosquatting, language plays a significant role in manipulating legitimate terms and names to create deceptive variations that appear familiar to the target audience. Attackers use linguistic techniques to construct these variations in ways that exploit the visual similarity or familiarity of certain languages and alphabets.



|ID| NAME        | GLYPHS      |HOMOPHONES  |ANTONYMS      |TYPOS  |CARDINAL  |ORDINAL  |STEMS |
|-----------|-------------|-------------|------------|--------------|-------|----------|---------|------|                                                         
|hy| Armenian |    38      |    1     |   1  |   1  |     24   |    0  |   0|
|fi| Finnish  |    29      |    1     |   1  |   1   |    11   |    1  |   0|
|fr| French    |   27      |    1      |  1  |   1   |    11  |    10   |  0|
|iw| Hebrew    |   22      |    2      |  1  |   5   |    11  |     0   |  0|
|fa| Persian   |   28      |    1      |  1  |   1    |   11  |     0  |   0|
|ru| Russian   |   41     |     1     |   1  |   1    |   44  |    10  |   0|
|ar |Arabic    |   28     |     1     |   1  |   0    |   11  |    11  |   0|
|en| English   |   26     |   485     |  93 | 4256    |   10  |     9  |   0|
|es| Spanish   |   27     |     1     |   1  |   1    |   31 |      4  |   0|



### Keyboard Layouts

Keyboard layouts are central to typosquatting because certain common typing errors are directly related to the physical arrangement of keys. For example, in the QWERTY layout, adjacent letters (like "e" and "r" or "i" and "o") are frequently mistyped, leading to common typos. Urlinsane can exploit these patterns by generating alternative domain names that reflect plausible mistakes, relying on users inadvertently typing close-but-incorrect keys.

Additionally, different keyboard layouts (such as AZERTY, QWERTZ, or Dvorak) produce unique typo patterns, allowing cybercriminals to target users in specific regions or with particular typing habits. By understanding these layout-specific errors, typosquatters can increase the likelihood of catching misdirected traffic, making keyboard layouts a significant factor in effective typosquatting strategies.


|  Arabic | Armenian  | English  | Finnish |  French   | Russian | Spanish | Hebrew  | Persian | 
|----------|----------|----------|---------|-----------|---------|--------|----------|---------|
| غفقثصض   | QWERTY   |  QWERTY  | QWERTY  | ACNOR     | ЯШЕРТЫ  | QWERTY | Standard | Farsi   |
| AZERTY   | QWERTY   |  AZERTY  |         |           | ЙЦУКЕН  | QWERTY |         |   |
| غفقثصض   |          |  QWERTZ  |         |           | ЙЦУКЕН  |        |         |  |
| QWERTY   |          |  DVORAK  |         |           |         |        |         ||




## Algorithms

Algorithms systematically generate plausible misspelled domain variations by analyzing common typing errors and linguistic patterns. These algorithms account for mistakes like adjacent key errors, omitted letters, and character swaps to create likely typo-based domains. More advanced algorithms leverage multi-lingual datasets, enabling the detection of typographical errors across different languages and keyboard layouts. This approach broadens the scope of potential typos, increasing protection against international typosquatting attempts.

| ID | Name                         | Description |
|----|------------------------------|-------------|
| di | Dot Insertion   | Inserting periods in the target domain name  |
| do | Dot Omission   | Omission periods in the target name   |
| dh | Dot/Hyphen Substitution  |  Swapping Dot and hyphen in a domain |
| hi | Hyphen Insertion | Inserting hyphens in the target domain            |
| ho | Hyphen Omission  | Created by removing hyphens from the domain|
| co | Character Omission | Omitting a character from the name            |
| cs | Character Substitution   |  Swapping two consecutive characters in a name |
|acs | Adjacent Char Sub | Replaces adjacent character from the keyboard    |
|aci | Adjacent Char Ins    |Inserting adjacent character from the keyboard |
|gi  | Grapheme Insertion |Inserting the language-specific alphabet in the target domain |
|gr  | Grapheme Replacement | Replaces an alphabet in the target domain|
| hr | Homoglyphs Replacement | Replaces characters with characters that look similar|
|sps | Singular Pluralise  | Singular forms of words are swapped for plural forms |
| cr | Character Repeat | Character Repeat Repeats a character of the domain name twice |
| dcr| Double Char Replacement    | Created by replacing identical, consecutive letters in the name.            |
|dcar| Double Char Adjacent Repl    | Double character adjacent replacement typos replace consecutive identical letters with adjacent keys on the keyboard  |
|cm  | Common Misspellings             | Created from a dictionary of commonly misspelled words in each language            |
| hs | Homophones Substitution |Substitutes words that sound the same but have different spellings |
|vs  | Vowel Substitution  | Vowel Swapping replaces vowels in the domain name with other vowels to create variations |
|bf  | Bitsquatting | Relies on random bit-errors to redirect connections |
|tld | Wrong TLD    | Wrong top level domain (TLD)            | 
|tld2| Wrong SLD | Wrong second level domain (TLD2) | 
|tld3| Wrong TLD3  | Wrong third level domain (TLD3) |
|ons | Ordinal Number Sub | Substituting digital numbers and ordinal numbers |
|cns | Cardinal Number Sub | Substituting digial numbers and carninal numbers  | 
| si | Subdomain Insertion | Inserts common subdomain at the beginning of the domain |
|com | Combosquatting | **TODO**: Combining keywords extracted via NLP and HTML meta tags |
| st | [Stem](https://en.wikipedia.org/wiki/Stemming) Substitution       | **TODO**: Substituting words with there root form |
| ks | Keyboard Substitution       | **TODO**: Changing international keyboard layouts while assuming the user is typing in his/her native layout.           |





## Information

Information gathering on target domains enables a detailed comparison of similar-looking domains to determine if they are being typosquatted by cybercriminals. By collecting data on domain ownership, registration dates, hosting locations, and site content, algorithms can analyze whether these variations are likely to be malicious. This approach helps identify suspicious patterns and potential connections to phishing, fraud, or brand impersonation attempts. With thorough data collection, organizations can better detect and respond to typosquatting threats in real time.



| ID  | Name              | Description                                                                                                    |
|-----|-------------------|------------------------------------------------------------------------------------------------|
|     | [Levenshtein](https://en.wikipedia.org/wiki/Levenshtein_distance) | Calculates Levenshtein distance between domains by default to limit scan scope.                   |
| a   | DNS A             | Retrieves host IPv4 addresses.                                                                      |
| mx  | DNS MX            | Retrieves DNS Mail Exchange (MX) records.                                                           |
| txt | DNS TXT           | Retrieves DNS TXT records.                                                                         |
| aa  | DNS AAAA          | Retrieves host IPv6 addresses.                                                                     |
| cn  | DNS CName         | Maps one domain to another via CNAME records.                                                      |
| ns  | DNS NS            | Checks NS records to identify the authoritative name server for a domain.                          |
| geo | GeoIP Info        | Provides IP location information via MaxMind database.                                             |
| ssd | SSDeep            | Uses fuzzy hashing with ssdeep to determine domain similarity, for pages with substantial content. |
| 301 | Redirects         | Retrieves domain redirects.                                                                        |
| idn | IDN               | Retrieves internationalized domain names.                                                          |
| bn  | Banner            | Captures HTTP/SMTP banner using a basic TCP connection.                                            |
| png | Screenshot        | Takes a domain screenshot via a headless browser and stores it locally.                            |
| wi  | Whois             | **TODO**: Perform Whois lookup for domain information.                                             |
| kw  | Keywords          | **TODO**: Extract keywords using the [RAKE](https://www.mathworks.com/help/textanalytics/ug/extract-keywords-from-documents-using-rake.html) algorithm. |
| tp  | NLP Topics        | **TODO**: Extract topics using the [LDA](https://en.wikipedia.org/wiki/Latent_Dirichlet_allocation) algorithm. |
| vc  | [VSM](https://en.wikipedia.org/wiki/Vector_space_model) | **TODO**: Compare domains' vector spaces for cosine similarity.                                    |
| lm  | [LLM](https://en.wikipedia.org/wiki/Large_language_model) | **TODO**: Use LLMs for keyword extraction, stemming, named entity recognition, and other NLP tasks. |
| ng  | [N-Gram](https://en.wikipedia.org/wiki/N-gram) | **TODO**: Generate domain variants using the domain's most common N-grams.                        |
| har | [HAR](https://en.wikipedia.org/wiki/HAR_(file_format)) | **TODO**: Retrieve HAR file from browser interaction for in-depth data analysis.

## Outputs

With structured outputs, users can seamlessly incorporate findings into their existing defenses, strengthening their protection against typosquatting threats.


| Name  | Description                               |
|-------|-------------------------------------------|
| TABLE | Pretty table format with color styling    |
| HTML  | HTML-formatted output                     |
| JSON  | **TODO**: JSON output format              |
| TXT   | Plain text output, one record per line    |
| CSV   | Comma-separated values format             |
| TSV   | Tab-separated values format               |
| MD    | Markdown-formatted output                 |

## In Progress

- I’m working on creating a `.urlinsane` directory in the user's home directory to store screenshots, data dumps, and cache, along with a configuration file to provide additional tool customization options.


## TODO

- [LLM](https://en.wikipedia.org/wiki/Large_language_model): I’m interested in exploring the use of Large Language Models (LLMs) to replace our individual natural language processing (NLP) algorithms and to automatically generate language datasets.

- I’d like to explore ways to reduce the program’s size, which is currently 11MB. Reusing existing datasets from the operating system, such as the MaxMind GeoIP database, TLD suffix lists, LLMs, and vector databases, could help minimize storage requirements.

- I’d like to explore restructuring the information-gathering functions using a DAG execution pattern with dependencies, rather than chaining plugins in a linear pipeline. This approach, similar to how Terraform handles plugin execution, would enable more flexible and efficient processing of interdependent tasks.

- I’d like to explore adding an analysis plugin that compares data between two domains and can be run as a separate CLI command.

- Create a script to download and build keyboard layouts from [kbdlayout.info](http://kbdlayout.info/)

- Create a more advanced model for a keyboard that includes layer-shifting 





###  Other Tools









## Authors

* [Rangertaha (rangertaha@gmail.com)](https://github.com/rangertaha)

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details






