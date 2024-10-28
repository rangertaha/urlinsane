# URLInsane

[![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)

Urlinsane is used to aid in the detection of typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage, and threat intelligence.

This is the most advanced and full-featured typosquatting tool.  It supports more algorithms(24), information gathering(8), keyboard layouts(19), and written languages(9) than any other typosquatting tool. It generates more results per domain in less time, see [`Speed`](#Speed). Originally, I wrote similar tools in Python but was not happy with the performance. This tool was built in Golang to take advantage of its speed, concurrency, and portability.  It builds around linguistic modeling, natural language processing (NLP), concurrency, and plugin architecture. It's easily extensible with plugins for typo-generation algorithms, information gathering, analysis, languages, keyboard layouts, and output formats. 


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
$ urlinsane 

Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking, 
URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.

Usage:
  urlinsane [flags]
  urlinsane [command]

Available Commands:
  help        Help about any command
  server      Start a websocket server to use this tool programmatically
  typo        Generates domain typos and variations

Flags:
      --config string   Configuration file (default is $HOME/.urlinsane.yaml)
  -h, --help            help for urlinsane

Use "urlinsane [command] --help" for more information about a command.

```

### Squatting Options

```bash
$ urlinsane typo -h


Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking,
URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.

USAGE:
  urlinsane typo [domains] [flags]

OPTIONS:
  -c, --concurrency int         Number of concurrent workers (default 50)
      --delay int               A delay between network calls (default 10)
  -f, --file string             Output filename
  -r, --filters stringArray     Filter results to reduce the number of results
  -o, --format string           Output format (csv, text) (default "text")
  -x, --funcs stringArray       Extra functions or filters (default [ld,idna])
  -h, --help                    help for typo
  -k, --keyboards stringArray   Keyboards/layouts ID to use (default [en])
      --random-delay int        Used to randomize the delay between network calls. (default 5)
  -t, --typos stringArray       Types of typos to perform (default [all])
  -v, --verbose                 Output additional details

GLOBAL OPTIONS:
      --config string   Configuration file (default is $HOME/.urlinsane.yaml)

TYPOS: 
  These are the types of typo/error algorithms that generate the domain variants
    MD	Missing Dot is created by omitting a dot from the domain.
    SI	Inserting common subdomain at the beginning of the domain.
    MDS	Missing Dashes is created by stripping all dashes from the domain.
    CO	Character Omission Omitting a character from the domain.
    CS	Character Swap Swapping two consecutive characters in a domain
    ACS	Adjacent Character Substitution replaces adjacent characters
    ACI	Adjacent Character Insertion inserts adjacent character 
    CR	Character Repeat Repeats a character of the domain name twice
    DCR	Double Character Replacement repeats a character twice.
    SD	Strip Dashes is created by omitting a dash from the domain
    SP	Singular Pluralise creates a singular domain plural and vice versa
    CM	Common Misspellings are created from a dictionary of commonly misspelled words
    VS	Vowel Swapping is created by swaps vowels
    HG	Homoglyphs replaces characters with characters that look similar
    WTLD	Wrong Top Level Domain
    W2TLD	Wrong Second Level Domain
    W3TLD	Wrong Third Level Domain
    HP	Homophones Modules are created from sets of words that sound the same
    BF	Bitsquatting relies on random bit-errors to redirect connections
    NS	Numeral Swap numbers, words and vice versa
    PI	Inserting periods in the target domain
    HI	Inserting hyphens in the target domain
    AI	Inserting the language specific alphabet in the target domain
    AR	Replacing the language specific alphabet in the target domain
    ALL	Apply all typosquatting algorithms

INFORMATION: 
  Retrieve aditional information on each domain variant.
    LD    The Levenshtein distance between strings
    IDNA    Show international domain name
    IP    Checking for IP address
    HTTP    Get http related information
    MX    Checking for DNS's MX records
    TXT    Checking for DNS's TXT records
    NS    Checks DNS NS records
    CNAME    Checks DNS CNAME records
    SIM    Show domain content similarity
    GEO    Show country location of IP address
    ALL    Apply all post typosquatting functions

FILTERS: 
  Filters to reduce the number of domain variants returned.
    LIVE   Show online/live domains only.
    ALL    Apply all filters

KEYBOARDS:
    AR1	Arabic keyboard layout
    AR2	Arabic PC keyboard layout
    AR3	Arabic North african keyboard layout
    AR4	Arabic keyboard layout
    HY1	Armenian QWERTY keyboard layout
    HY2	Armenian, Western QWERTY keyboard layout
    EN1	English QWERTY keyboard layout
    EN2	English AZERTY keyboard layout
    EN3	English QWERTZ keyboard layout
    EN4	English DVORAK keyboard layout
    FI1	Finnish QWERTY keybaord layout
    FR1	French Canadian CSA keyboard layout
    IW1	Hebrew standard layout
    FA1	Persian standard layout
    RU1	Russian keyboard layout
    RU2	Phonetic Russian keybaord layout
    RU3	PC Russian keyboard layout
    ES1	Spanish keyboard layout
    ES2	Spanish ISO keyboard layout
    ALL	Use all keyboards

EXAMPLE:

    urlinsane google.com
    urlinsane google.com -t co
    urlinsane google.com -t co -x ip -x idna -x ns

AUTHOR:
    Written by Rangertaha


```

## Server Options

```bash

urlinsane server -h

Usage:
  urlinsane server [flags]

Flags:
  -c, --concurrency int   Number of concurrent workers (default 50)
  -h, --help              help for server
  -a, --host string       IP address for API server (default "127.0.0.1")
  -p, --port string       Port to use (default "8080")

Global Flags:
      --config string   Configuration file (default is $HOME/.urlinsane.yaml)

```

## Usage

Generates variations for **google.com** using the character omission **(CO)**
algorithm.

```txt
urlinsane typo google.com -t co

 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: 0.6.0

  LIVE | TYPE |   TYPO    | SUFFIX |   IDNA
-------+------+-----------+--------+------------
       | CO   | oogle.com | com    | oogle.com  
       | CO   | gogle.com | com    | gogle.com  
       | CO   | goole.com | com    | goole.com  
       | CO   | gogle.com | com    | gogle.com  
       | CO   | googl.com | com    | googl.com  
       | CO   | googe.com | com    | googe.com  

```

Additional e**x**tra functions can be selected with the **-x, --funcs** options.
These functions can add columns to the output. For example, the following generates
variations for **google.com** using the character omission **(CO)** algorithm
then checks for **ip** addresses.

```txt

urlinsane typo google.com -t co  -x geo

```

Generates variations for **google.com** with the following parameters:

* **-t hg** lets us use the Homoglyphs(HG) algorithm only
* **-v** Verbose mode shows us the full name 'Homoglyphs' of the algorithm not
just the short name 'HG'
* **-x ip** Check or IP address
* **-x idna** Shows the IDNA format
* **-x ns** Checks for DNS NS records

```txt

urlinsane typo google.com -t hg -v -x ip -x idna -x ns


```

## Languages

### English

* Over 8000 common misspellings
* Over 500 common homophones
* English alphabet, vowels, homoglyphs, and numerals
* Common keyboard layouts (qwerty, azerty, qwertz, dvorak)

### Finnish, Russian, Persian, Hebrew, Arabic, Spanish

See [Languages](https://cybersectech-org.github.io/urlinsane/#languages) for details
on other languages.

## Algorithms

The modular architecture for code extensibility allows developers to add new
typosquatting algorithms with ease. Currently, we have implemented 19
typosquatting algorithms. See [Typo Algorithms](https://cybersectech-org.github.io/urlinsane/#algorithms) for details.

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

## Post Typo Functions

|      **Extra Functions**            | URLInsane  | URLCrazy  | DNSTwist  | DomainFuzz  |
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






