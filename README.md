# URLInsane

[![Build Status](https://travis-ci.org/cybersectech-org/urlinsane.svg?branch=master)](https://travis-ci.org/cybersectech-org/urlinsane) 

Multilingual domain typo permutation engine used to perform or detect typosquatting, 
brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage and 
threat intelligence.

**Documentation:** [URLInsane Docs](https://cybersectech-org.github.io/urlinsane/)

**Downloads:** [URLInsane Downloads](https://github.com/cybersectech-org/urlinsane/releases/tag/0.3.0)


<div align="center">
  <a href="https://www.youtube.com/watch?v=_IU1ADTBKVU" style="float:left"><img width="350" src="https://img.youtube.com/vi/_IU1ADTBKVU/0.jpg" alt="URLInsane Demo1"></a>
   <a href="https://www.youtube.com/watch?v=HgMV0NqMCm0" style="float:left"><img width="350"  src="https://img.youtube.com/vi/HgMV0NqMCm0/0.jpg" alt="URLInsane Demo2"></a>
</div>




## Features

* Binary executable, written in Go with no dependencies. 
* Will have all the functionally of URLCrazy and DNSTwist. 
* Contains 22 typosquatting algorithms and 10 extra functions to retrieve additional data such as ip to geographic location, dns lookups and more 
* Modular architecture for language, keyboard, typo algorithm, and functions extensibility.
* Supports multiple keyboard layouts found in English, Spanish, Russian, Armenian, Finish, French, Hebrew, Persian, and Arabic.
* Supports multiple languages with the ability to add more languages with ease.
* Concurrent function (**-x --funcs**) workers to retrieve additional info on each record.
* Concurrent typo squatting workers.

## Cli Tool

```
Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.

Usage:
  urlinsane [command]

Available Commands:
  help        Help about any command
  server      Start an API server to use this tool programmatically
  typo        Generates domain typos and variations

Flags:
  -h, --help   help for urlinsane

Use "urlinsane [command] --help" for more information about a command.
```

### Squatting Options

```
urlinsane typo -h


Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.

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

TYPOS: These are the types of typo/error algorithms that generate the domain variants
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
  ALL   Apply all typosquatting algorithms

FUNCTIONS: Post processig functions that retieve aditional information on each domain variant.
  LD	The Levenshtein distance is a string metric for measuring the difference between two domains
  IDNA	Show international domain name
  IP	Checking for IP address
  HTTP	Get http related information
  MX	Checking for DNS's MX records
  TXT	Checking for DNS's TXT records
  NS	Checks DNS NS records
  CNAME	Checks DNS CNAME records
  SIM	Show domain content similarity
  GEO	Show country location of ip address
  ALL  	Apply all post typosquating functions

FILTERS: Filters to reduce the number domain variants returned.
  LIVE	Show online/live domains only.
  ALL  	Apply all filters

EXAMPLE:

    urlinsane google.com
    urlinsane google.com -t co
    urlinsane google.com -t co -x ip -x idna -x ns

AUTHOR:
  Written by Rangertaha <rangertaha@gmail.com>
```

## Server Options
```
urlinsane server -h

This command starts up a REST API server to use this tool programmatically.

Usage:
  urlinsane server [flags]

Flags:
  -a, --addr.host string   IP address for API server (default "127.0.0.1")
  -p, --addr.port string   Port to use (default "8888")
  -c, --concurrency int    Number of concurrent workers (default 50)
  -h, --help               help for server
      --log.file string    Filename to send logs to (default "urlinsane.log")
      --log.level string   Logging level (default "DEBUG")
  -s, --stream             Stream results via http2
```

## Usage
Generates variations for **google.com** using the character omission **(CO)** 
algorithm.

```
urlinsane typo google.com -t co

 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: 0.5.2

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
These functions can add columns to the output. For example the following generates 
variations for **google.com** using the character omission **(CO)** algorithm 
then checks for **ip** addresses. 

```
urlinsane typo google.com -t co  -x geo

```

Generates variations for **google.com** with the following parameters:
* **-t hg** lets us use the Homoglyphs(HG) algorithm only
* **-v** Verbose mode shows us the full name 'Homoglyphs' of the algorithm not 
just the short name 'HG'
* **-x ip** Check or ip address
* **-x idna** Shows the IDNA format
* **-x ns** Checks for DNS NS records

```
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
typosquatting algorithms with ease. Currently we have implements 19 
typosquatting algorithms. See [Typo Algorithms](https://cybersectech-org.github.io/urlinsane/#algorithms) for details.


## Extra Functions

- **IDNA**  Show international domain name (Default)
- **MX**    Checking for DNS's MX records
- **TXT**   Checking for DNS's TXT records
- **IP**    Checking for IP address
- **NS**    Checks DNS NS records
- **CNAME** Checks DNS CNAME records
- **SIM**   Show domain similarity % using fuzzy hashing with ssdeep
- **LIVE**	Show domains with ip addresses only
- **301**	Show domains redirects
- **GEO**	Show country location of ip address



## Tools Comparisons

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


### Algorithms

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
| Combosquatting(Keywords)        |           |           |            |             |           
| Multithreaded Algorithms        |     X     |     ?     |      X     |             |
| Subdomain insertion             |     X     |           |            |             |
| Period Insertion                |     X     |           |            |             |  

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
| GeoIP Info                          |     X     |     X     |      X     |             |           
| Domain Similarity                   |     X     |           |      X     |      X      |           
| Domain Redirects                    |     X     |           |            |             |           
| IDNA Format                         |     X     |           |      X     |             |           
| CSV output                          |     X     |     X     |      X     |      X      |           
| JSON output                         |     X     |           |      X     |      X      |           
| Human Readable output               |     X     |     X     |      X     |      X      | 
| Google Popularity Estimate          |           |     X     |            |             |           
| HTTP/SMTP Banner                    |     X     |           |      X     |             |           
| WHOIS Info                          |           |           |      X     |             |                    
| Multithreaded Extra Functions       |     X     |           |      X     |      X      |           







### TODO 

* Extract keywords from domains. Keywords will be used for additional algorithms
* Estimate popularity of a domain variant via google search
* Lookup whois record
* Emoji domains
* Grabs HTTP and SMTP service banners
* Speed Comparison with other tools


## Authors

* [Rangertaha <rangertaha@gmail.com>](https://github.com/rangertaha)


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
