# URLInsane

[![Go Report Card](https://goreportcard.com/badge/github.com/rangertaha/urlinsane?style=flat-square)](https://goreportcard.com/report/github.com/rangertaha/urlinsane) [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/rangertaha/urlinsane) [![PkgGoDev](https://pkg.go.dev/badge/github.com/rangertaha/urlinsane)](https://pkg.go.dev/github.com/github.com/rangertaha/urlinsane) [![Release](https://img.shields.io/github/release/rangertaha/urlinsane.svg?style=flat-square)](https://github.com/rangertaha/urlinsane/releases/latest) [![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)

Urlinsane is used to aid in the detection of typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage, supply chain attacks, and threat intelligence. It's a command-line tool for detecting typosquatting across domains, usernames, and software packages. It scans for potential typosquatting variants by applying advanced typo squatting algorithms, information gathering, and data analysis.  It identifies potentially harmful variations of a victim's domain name, arbitrary names, and software packages that cybercriminals might exploit. 

It's inspired by [URLCrazy](https://morningstarsecurity.com/research/urlcrazy), [Dnstwist](https://github.com/elceef/dnstwist), [DomainFuzz](https://github.com/monkeym4ster/DomainFuzz) and a few other libraries and tools I was researching at the time.


## Targets

 Urlinsane supports typo generation and information collection for **Domains**, **Emails**, **Usernames**, and software **Packages**.
## Domain 

```bash
urlinsane typo -d example.com 
```

## Email 
```bash
urlinsane typo -e username@example.com 
```

## Username 

```bash
urlinsane typo -n urlinsane
```

```bash
urlinsane typo -n urlinsane -u https://github.com/rangertaha/urlinsane
```

## Packages 

```bash
urlinsane typo -g express
```


```bash
urlinsane typo -g express -u https://www.npmjs.com/package/express
```



## Installation


## Usage




# Internals


## Plugins

Plugins play a crucial role in extending the functionality, flexibility, and customization of Urlinsane and allow it to evolve alongside changing needs and technological advancements. Here's a structured summary of the plugin types and their roles in Urlinsane:

|    Type     | Number  | Description |
|-------------|--------|--------------|
| Languages   |    9   | Language plugins provide data for it's linguistic capability. |
| Keyboards   |    19  | Keyboard plugins provide layouts for international keyboads |
| Algorithms  |    24  | They generate typo variants for each target type |
| Information |    8   | Collects information on target types |
| Database    |    1   | TODO: Caches and saves results for analysis |
| Outputs     |    6   | Formats and saves results  |


### Languages

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

### Keyboard Layouts

| Arabic | Armenian | English  | Finnish | Russian   | Spanish | Hebrew | Persian | 
|----------|----------|----------|---------|-----------|---------|--------|---------|
| "غفقثصض" |          |  QWERTY  |         |           |         |        |         |
|        |          |  AZERTY  |         |           |         |        |         |   
|        |          |  QWERTZ  |         |           |         |        |         |  
|        |          |  DVORAK  |         |           |         |        |         |
|        |          |          |         |           |         |        |         | 
|        |          |          |         |           |         |        |         |  
|        |          |          |         |           |         |        |         | 
|        |          |          |         |           |         |        |         | 


## Algorithms


| ID | Name                         | Description |
|----|------------------------------|-------------|
|    | Missing Dot                  |             |
|    | Missing Dashes               |             |
|    | Strip Dashes                 |             |
|    | Character Omission           |             |
|    | Character Swap                  |             |
|    | Adjacent Character Substitution |             |
|    | Adjacent Character Insertion    |             |
|    | Homoglyphs                      |             |
|    | Singular Pluralise              |             |
|    | Character Repeat                |             |
|    | Double Character Replacement    |             |
|    | Common Misspellings             |             |
|    | Homophones                      |             |
|    | Vowel Swapping                  |             |
|    | Bitsquatting                    |             |
|    | Wrong Top Level Domain          |             | 
|    | Wrong Second Level Domain       |             | 
|    | Wrong Third Level Domain        |             |
|    | Ordinal Number Swap             |             |
|    | Cardinal Number Swap            |             |
|    | Hyphenation                     |             | 
|    | Multithreaded Algorithms        |             |   
|    | Subdomain insertion             |             |
|    | Period Insertion                |             | 
|    | Combosquatting (Keywords)       |             |



## Information

### Domain Information

| ID |  Name             | Description  |
|----|-------------------|--------------|
|    | DNS A Records     | Retrieving IPv4 and IPv6 IP host addresses |
|    | DNS MX Records    | Retrieving Mail Exchange (MX) records|
|    | DNS TXT Records   | Retrieving TXT records storing arbitrary data associated with a domain |
|    | DNS AAAA Records  | |
|    | DNS CName Records | |
|    | DNS NS Records    | Checks DNS NS records |
|    | Geographic Info   | Show country location of IP address|
|    | Domain Similarity | Show domain similarity % using fuzzy hashing with ssdeep|
|    | Domain Redirects  | Show domains redirects |
|    | IDNA Format       | Show international domain name (Default) |
|    | HTTP/SMTP Banner  | |

### Package Information

| ID |  Name                | Description  |
|----|----------------------|--------------|
| py | Python Package Index | Checks PyPi servers for pacakges |
| js | Node Package Manager | Checks NPM servers for pacakges |
| gh | Github               | Checks NPM servers for pacakges |



### Username Information

| ID |  Name             | Description  |
|----|-------------------|--------------|
|    |                   |              |





## Outputs

| Name  | Description | 
|-------|-------------|
| TABLE |         |  
| TEXT  |         |  
| CSV   |         |    
| TSV   |         |   
| MD    |         |   


## Database

| Name   | Description | 
|--------|-------------|
| Badger |             |    







###  Other Tools









## Authors

* [Rangertaha (rangertaha@gmail.com)](https://github.com/rangertaha)

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details






