# URLInsane

[![Go Report Card](https://goreportcard.com/badge/github.com/rangertaha/urlinsane?style=flat-square)](https://goreportcard.com/report/github.com/rangertaha/urlinsane) [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/rangertaha/urlinsane) [![PkgGoDev](https://pkg.go.dev/badge/github.com/rangertaha/urlinsane)](https://pkg.go.dev/github.com/github.com/rangertaha/urlinsane) [![Release](https://img.shields.io/github/release/rangertaha/urlinsane.svg?style=flat-square)](https://github.com/rangertaha/urlinsane/releases/latest) [![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)

Urlinsane is used to aid in the detection of typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage, supply chain attacks, and threat intelligence. It's a command-line tool for detecting typosquatting domains. It scans for potential typosquatting variants by applying advanced typo squatting algorithms, information gathering, and data analysis.  It identifies potentially harmful variations of a victim's domain name that cybercriminals might exploit. 

It's inspired by [URLCrazy](https://morningstarsecurity.com/research/urlcrazy), [Dnstwist](https://github.com/elceef/dnstwist), [DomainFuzz](https://github.com/monkeym4ster/DomainFuzz) and a few other libraries and tools I was researching at the time.





## Installation


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
| Algorithms  |    24  | They generate typo variants for each target type |
| Information |    13  | Collects information on target types |
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

|  Arabic | Armenian  | English  | Finnish |  French   | Russian | Spanish | Hebrew  | Persian | 
|----------|----------|----------|---------|-----------|---------|--------|----------|---------|
| غفقثصض   | QWERTY   |  QWERTY  | QWERTY  | ACNOR     | ЯШЕРТЫ  | QWERTY | Standard | Farsi   |
| AZERTY   | QWERTY   |  AZERTY  |         |           | ЙЦУКЕН  | QWERTY |         |   |
| غفقثصض   |          |  QWERTZ  |         |           | ЙЦУКЕН  |        |         |  |
| QWERTY   |          |  DVORAK  |         |           |         |        |         ||



## Algorithms


| ID | Name                         | Description |
|----|------------------------------|-------------|
| di | Dot Insertion                 |             |
| do | Dot Omission                 |             |
| dh | Dot/Hyphen Substitution                 |             |
| hi | Hyphen Insertion               |             |
| ho | Hyphen Omission                |             |
| co | Character Omission           |             |
| cs | Character Substitution                  |             |
|acs | Adjacent Char Sub |             |
|aci | Adjacent Char Ins    |             |
|gi  | Grapheme Insertion                      |             |
|gr  | Grapheme Replacement                      |             |
| hr | Homoglyphs Replacement                      |             |
|sps | Singular Pluralise              |             |
| cr | Character Repeat                |             |
| dcr| Double Char Replacement    |             |
|dcar| Double Char Adjacent Repl    |             |
|cs  | Common Misspellings             |             |
| hs | Homophones Substitution                     |             |
|vs  | Vowel Substitution                  |             |
|bf  | Bitsquatting                    |             |
|tld | Wrong Top Level Domain          |             | 
|tld2| Wrong Second Level Domain       |             | 
|tld3| Wrong Third Level Domain        |             |
|ons | Ordinal Number Sub             |             |
|cns | Cardinal Number Sub           |             | 
| si | Subdomain insertion             |             |
|com | Combosquatting        | **TODO**: Combining keywords extracted via NLP and HTML meta tags |
| st | [Stem](https://en.wikipedia.org/wiki/Stemming) Substitution       | TODO            |
| ks | Keyboard Substitution       | **TODO**: Changing international keyboard layouts while assuming the user is typing in his/her native layout.           |


## Information


| ID  | Name       | Description  |
|-----|-------------------|--------------|
| ld   | [Levenshtein](https://en.wikipedia.org/wiki/Levenshtein_distance)     | Levenshtein distance bettween domains is created by defualt and used internally to limit the number of scans |
| a   | DNS A      | Host IPv4 addresses |
| mx  | DNS MX     | DNS Mail Exchange (MX) records|
| txt | DNS TXT    | DNS TXT records records |
| aa  | DNS AAAA   | Host IPv6 addresses |
| cn  | DNS CName  | CName records are used to maps one domain to another |
| ns  | DNS NS     | Checks NS records specifying the authoritative name server for a domain |
| geo | GeoIP Info | Show country location of IP address via MaxMind database|
| ssd | SSDeep     | Domain similarity using fuzzy hashing with ssdeep, only works with pages with enough content|
| 301 | Redirects  | Get domains redirects |
| idn | IDN        | Get international domain names |
| bn  | Banner     | HTTP/SMTP Banner using a simple TCP connection |
| png | Screenshot | Screenshot of the domain via headless browser and stores locally |
| wi  | Whois      | **TODO**: Whois domain loookup info |
| kw  | Keywords   | **TODO**: Extracting keywords use the [RAKE](https://www.mathworks.com/help/textanalytics/ug/extract-keywords-from-documents-using-rake.html) algorithm |
| tp  | NLP Topics | **TODO**: Extracting topics via [LDA](https://en.wikipedia.org/wiki/Latent_Dirichlet_allocation) algorithm |
| vc  | [VSM](https://en.wikipedia.org/wiki/Vector_space_model)    | **TODO**: Comparing domains vector space for cosin similarity |
| lm  | [LLM]()    | **TODO**: Using LLMs for keywork extraction, stemming, named entity extraction and a lot more NLP functionality |


## Outputs

| Name  | Description | 
|-------|-------------|
| TABLE |Pretty table output format with color         |  
| HTML  | HTML formatted output        |   
| JSON  | **TODO**: JSON outputs format        |  
| TXT   | Text outputs streams one record per line        | 
| CSV   | CSV (comma-separated values) formatted output        |    
| TSV   | TSV (tab-separated values) formatted output        |   
| MD    | Markdown formatted output         |   






###  Other Tools









## Authors

* [Rangertaha (rangertaha@gmail.com)](https://github.com/rangertaha)

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details






