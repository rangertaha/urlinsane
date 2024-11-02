# URLInsane

[![Go Report Card](https://goreportcard.com/badge/github.com/rangertaha/urlinsane?style=flat-square)](https://goreportcard.com/report/github.com/rangertaha/urlinsane) [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/rangertaha/urlinsane) [![PkgGoDev](https://pkg.go.dev/badge/github.com/rangertaha/urlinsane)](https://pkg.go.dev/github.com/github.com/rangertaha/urlinsane) [![Release](https://img.shields.io/github/release/rangertaha/urlinsane.svg?style=flat-square)](https://github.com/rangertaha/urlinsane/releases/latest) [![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)


URLInsane is a robust command-line tool designed to detect typosquatting across domains, arbitrary names, usernames, and software packages. By leveraging advanced algorithms, information-gathering techniques, and data analysis, it identifies potentially harmful variations of targeted entities that cybercriminals might exploit. Essential for defending against typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, and corporate espionage, URLInsane also enhances threat intelligence capabilities.

Featuring a plugin-based multilingual permutation engine, URLInsane supports various keyboard layouts for multiple languages. Its extensible plugin system allows easy addition of new capabilities. Currently, it includes plugins for 24 algorithms, 8 information gathering methods, 19 keyboard layouts, 9 languages, and 4 output formats. Originally developed in Python, URLInsane was built in Go to enhance speed, concurrency, and portability.



Urlinsane is a powerful command-line tool designed for the detection of typo-squatting across domains, email addresses, usernames, and software packages. Utilizing  Urlinsane scans for potential typosquatted variants by applying advanced typo squatting algorithms, information gathering, and data analysis.  It identifies potentially harmful variations of a victim's domain name, email address, software packages, and that cybercriminals might exploit. 



Urlinsane is used to aid in the detection of typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage, and threat intelligence.

It's a plugin-based multilingual permutation engine that supports keyboard layouts for each language. The plugin system allows us to easily extend its capabilities. Currently, it supports plugins for algorithms (24), information gathering (8), keyboard layouts (19), languages (9), and output (4) formats. Originally, I wrote similar tools in Python but was not happy with the performance. This tool was built in Golang to take advantage of its speed, concurrency, and portability (see [`Speed`](#Speed).




| **Plugins** |  count |
|-------------|--------|
| Languages   |    9   |
| Keyboards   |    19  |
| Algorithms  |    24  |
| Information |    8   |
| Analysis    |    2   |



 


## Features

* Binary executable, written in Go with no dependencies.
* Will have all the functionality of URLCrazy and DNSTwist.
* Contains 24 typosquatting algorithms and 10 extra functions to retrieve additional data such as IP to geographic location, DNS lookups and more
* Modular architecture for language, keyboard, typo algorithm, and functions extensibility.
* Supports multiple keyboard layouts found in English, Spanish, Russian, Armenian, Finnish, French, Hebrew, Persian, and Arabic.
* Supports multiple languages with the ability to add more languages with ease.
* Concurrent function (**-x --funcs**) workers to retrieve additional info on each record.
* Concurrent typo squatting workers.



## Example

Finds "character omission" typos for the given domain. **-t** specifies the type of typo you want to use and defaults to 
all 24. **-x** specifies the extra information retrieval functions to use and defaults to non-internet required functions. 
 
```bash
$ urlinsane typo google.com -t co -x all 

 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: 0.7.0

   LIVE  | TYPE |   TYPO    | SUFFIX | LD |   IDNA    |      IPV4      |           IPV6           | SIZE |    REDIRECT    |        MX        |                                            TXT                                             |           NS           | CNAME | SIM |      GEO       
---------+------+-----------+--------+----+-----------+----------------+--------------------------+------+----------------+------------------+--------------------------------------------------------------------------------------------+------------------------+-------+-----+----------------
  ONLINE | CO   | googl.com | com    |  1 | googl.com | 172.217.10.228 | 2607:f8b0:4006:813::2004 |      | www.google.com |                  | v=spf1 -all                                                                                | ns3.google.com         |       |     | United States  
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns2.google.com         |       |     |                
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns4.google.com         |       |     |                
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns1.google.com         |       |     |                
  ONLINE | CO   | oogle.com | com    |  1 | oogle.com | 104.28.29.162  | 2606:4700:30::681c:1da2  |      |                | mx.zoho.com      | brave-ledger-verification=2dd5f8cc6d7ac0d6d6f27de1c07629a8e329ecdebdc7303506854fc3eec20968 | gwen.ns.cloudflare.com |       |     | United States  
         |      |           |        |    |           | 104.28.28.162  | 2606:4700:30::681c:1ca2  |      |                | mx2.zoho.com     | v=spf1 +a +mx +ip4:204.9.184.9 +include:zoho.com ~all                                      | amir.ns.cloudflare.com |       |     |                
  ONLINE | CO   | gogle.com | com    |  1 | gogle.com | 172.217.10.132 | 2607:f8b0:4006:810::2004 |      | www.google.com |                  | v=spf1 -all                                                                                | ns4.google.com         |       |     | United States  
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns2.google.com         |       |     |                
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns1.google.com         |       |     |                
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns3.google.com         |       |     |                
  ONLINE | CO   | goole.com | com    |  1 | goole.com | 217.160.0.201  |                          |      | www.goole.com  | mx00.1and1.co.uk |                                                                                            | ns1083.ui-dns.com      |       |     | Germany        
         |      |           |        |    |           |                |                          |      |                | mx01.1and1.co.uk |                                                                                            | ns1083.ui-dns.biz      |       |     |                
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns1083.ui-dns.de       |       |     |                
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns1083.ui-dns.org      |       |     |                
  ONLINE | CO   | googe.com | com    |  1 | googe.com | 50.63.202.32   |                          |      |                |                  | v=spf1 -all                                                                                | ns2.yourdoor.com       |       |     | United States  
         |      |           |        |    |           |                |                          |      |                |                  |                                                                                            | ns1.yourdoor.com       |       |     |                


```



## Cli Commands

```bash


```

### Squatting Options

```bash

```

## Server Options

```bash

```

## Usage

Generates variations for **google.com** using the character omission **(CO)**
algorithm.

```txt

```

Additional e**x**tra functions can be selected with the **-x, --funcs** options.
These functions can add columns to the output. For example, the following generates
variations for **google.com** using the character omission **(CO)** algorithm
then checks for **ip** addresses.

```txt


```

Generates variations for **google.com** with the following parameters:

* **-t hg** lets us use the Homoglyphs(HG) algorithm only
* **-v** Verbose mode shows us the full name 'Homoglyphs' of the algorithm not
just the short name 'HG'
* **-x ip** Check or IP address
* **-x idna** Shows the IDNA format
* **-x ns** Checks for DNS NS records

```txt


```

## Languages

### Language & Keyboard Comparison

This table shows which tools have support for common **misspellings**,
**homophones**, **numerals**, **vowels**, **homoglyphs**, and the number of
**keyboards** that support each language's character set.

Vowels
Graphemes
Ordinal
Cardinal

| Languages | Keyboards | Homoglyphs | Homophones | Antonyms | Misspellings | Vowels | Graphemes | Ordinal | Cardinal  | 
|-----------|-----------|-----------|------------|-----------|--------------|--------|-----------|---------|-----------|
| Arabic    |    4      |           |            |           |              |   |   |   |    | 
| Armenian  |    3      |           |            |           |              |   |   |   |    | 
| English   |    4      |           |            |           |              |   |   |   |    | 
| Finnish   |    1      |           |            |           |              |   |   |   |    | 
| Russian   |    3      |           |            |           |              |   |   |   |    | 
| Spanish   |    2      |           |            |           |              |   |   |   |    | 
| Hebrew    |    1      |           |            |           |              |   |   |   |    | 
| Persian   |    1      |           |            |           |                |   |   |   |    | 











### English

* Over 8000 common misspellings
* Over 500 common homophones
* English alphabet, vowels, homoglyphs, and numerals
* Common keyboard layouts (qwerty, azerty, qwertz, dvorak)

### Finnish, Russian, Persian, Hebrew, Arabic, Spanish

See [Languages](https://github.com/rangertaha/urlinsane#languages) for details
on other languages.

## Algorithms

The modular architecture for code extensibility allows developers to add new
typosquatting algorithms with ease. Currently, we have implemented 19
typosquatting algorithms. See [Typo Algorithms](https://github.com/rangertaha/urlinsane#algorithms) for details.

## Extra Functions

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

## Tools Comparisons

### Results

| **Tool**   | google.com  | facebook.com  | youtube.com   | amazon.com | amazon4you.com |
|------------|-------------|---------------|---------------|------------|----------------|
| URLInsane  |    6931     |    7049       |     6996      |   6934     |     7192       |
| URLCrazy   |    88       |    109        |     107       |   78       |     129        |
| DNSTwist   |    1687     |    2529       |     3770      |   2262     |     5815       |

### Language & Keyboard Comparison

This table shows which tools have support for common **misspellings**,
**homophones**, **numerals**, **vowels**, **homoglyphs**, and the number of
**keyboards** that support each language's character set.

| **Lang (# Keyboards)**   | URLInsane  | URLCrazy  | DNSTwist   | DomainFuzz |
|--------------------------|-----------|-----------|------------|-------------|
| Arabic (4)               |     X     |           |            |             |
| Armenian (3)             |     X     |           |            |             |
| English (4)              |     X     |     X     |      X     |      X      |
| Finnish (1)              |     X     |           |            |             |
| Russian (3)              |     X     |           |            |             |
| Spanish (2)              |     X     |           |            |             |
| Hebrew (1)               |     X     |           |            |             |
| Persian (1)              |     X     |           |            |             |  

### Algorithm Comparisons

This table shows the list of algorithms supported for each tool.

|      **Algorithms**             | URLInsane | URLCrazy  | DNSTwist   | DomainFuzz **(TODO)**  |
|---------------------------------|-----------|-----------|------------|-------------|
| Missing Dot                     |     X     |     X     |     X      |             |
| Missing Dashes                  |     X     |           |            |             |
| Strip Dashes                    |     X     |     X     |            |             |
| Character Omission              |     X     |     X     |     X      |             |
| Character Swap                  |     X     |     X     |            |             |
| Adjacent Character Substitution |     X     |     X     |            |             |
| Adjacent Character Insertion    |     X     |     X     |     X      |             |
| Homoglyphs                      |     X     |     X     |     P      |             |
| Singular Pluralise              |     X     |     X     |            |             |
| Character Repeat                |     X     |     X     |     X      |             |
| Double Character Replacement    |     X     |     X     |            |             |
| Common Misspellings             |     X     |     X     |            |             |
| Homophones                      |     X     |     X     |     P      |             |
| Vowel Swapping                  |     X     |     X     |            |             |
| Bitsquatting                    |     X     |     X     |     X      |             |
| Wrong Top Level Domain          |     X     |     X     |            |             |
| Wrong Second Level Domain       |     X     |     X     |            |             |
| Wrong Third Level Domain        |     X     |           |            |             |
| Ordinal Number Swap             |     X     |           |            |             |
| Cardinal Number Swap            |     X     |           |            |             |
| Hyphenation                     |     X     |           |      X     |             |
| Multithreaded Algorithms        |     X     |     ?     |      X     |             |
| Subdomain insertion             |     X     |           |            |             |
| Period Insertion                |     X     |           |            |             |
| Combosquatting (Keywords)       |           |           |            |             |

## Information Gathering

|      **Info Gathering**            | URLInsane  | URLCrazy  | DNSTwist  | DomainFuzz  |
|-------------------------------------|-----------|-----------|------------|-------------|
| Live/Online Check                   |     X     |     X     |      X     |             |
| DNS A Records                       |     X     |     X     |      X     |      X      |
| DNS MX Records                      |     X     |     X     |      X     |             |
| DNS txt Records                     |     X     |     X     |            |             |
| DNS AAAA Records                    |     X     |           |      X     |      X      |
| DNS CName Records                   |     X     |           |            |             |
| DNS NS Records                      |     X     |           |      X     |      X      |
| Geographic Info                     |     X     |     X     |      X     |             |
| Domain Similarity                   |     X     |           |      X     |      X      |
| Domain Redirects                    |     X     |           |            |             |
| IDNA Format                         |     X     |           |      X     |             |
| CSV output                          |     X     |     X     |      X     |      X      |
| JSON output                         |     X     |           |      X     |      X      |
| Human Readable output               |     X     |     X     |      X     |      X      |
| HTTP/SMTP Banner                    |     X     |           |      X     |             |
| Multithreaded Extra Functions       |     X     |           |      X     |      X      |

## Output Formats

|      **Info Gathering**            | URLInsane  | URLCrazy  | DNSTwist  | DomainFuzz  |
|-------------------------------------|-----------|-----------|------------|-------------|
| Live/Online Check                   |     X     |     X     |      X     |             |
| DNS A Records                       |     X     |     X     |      X     |      X      |
| DNS MX Records                      |     X     |     X     |      X     |             |
| DNS txt Records                     |     X     |     X     |            |             |
| DNS AAAA Records                    |     X     |           |      X     |      X      |
| DNS CName Records                   |     X     |           |            |             |
| DNS NS Records                      |     X     |           |      X     |      X      |
| Geographic Info                     |     X     |     X     |      X     |             |
| Domain Similarity                   |     X     |           |      X     |      X      |
| Domain Redirects                    |     X     |           |            |             |
| IDNA Format                         |     X     |           |      X     |             |
| CSV output                          |     X     |     X     |      X     |      X      |
| JSON output                         |     X     |           |      X     |      X      |
| Human Readable output               |     X     |     X     |      X     |      X      |
| HTTP/SMTP Banner                    |     X     |           |      X     |             |
| Multithreaded Extra Functions       |     X     |           |      X     |      X      |



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






