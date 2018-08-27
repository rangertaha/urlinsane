# URLInsane

Generates domain typos and variations used to detect and perform typo squatting, URL hijacking, phishing, and corporate espionage.
Inpired by URLCrazy I wanted to create a better version that supported multiple languages and linguistic typos.
I also wanted it to be a binary with fast execution time.


## Introduction
Generate and test domain typos and variations to detect and perform typo squatting, URL hijacking, phishing, and corporate espionage.


## Installation

Create the binary executable with the make command or [download](https://github.com/rangertaha/urlinsane/releases) one of the pre-built release binaries. 

```bash
make
```

Execute tool showing help
```bash
./urlinsane -h

    
Generates domain typos and variations to detect and perform typo squatting, URL hijacking, phishing, and corporate espionage.

USAGE:
  urlinsane [domains] [flags]

OPTIONS:
  -c, --concurrency int         Number of concurrent workers (default 50)
  -f, --file string             Output filename
  -o, --format string           Output format (json, csv, text) (default "text")
  -x, --funcs stringArray       Extra functions for retrieving additional data (default [idna])
  -h, --help                    help for urlinsane
  -k, --keyboards stringArray   Keyboards/layouts ID to use (default [en1])
  -t, --typos stringArray       Types of typos to perform (default [all])
  -v, --verbose                 Output additional details
  
  ...
  
  
```





## Features
* Binary executable 
* Multiple keyboard layouts
* Multiple languages




## Algorithms

URLInsane implements 19 typosquatting algorithms. 

MD    Missing Dot is created by omitting a dot from the domain.
MDS   Missing Dashes is created by omitting a dash from the domain.
SD    Strip Dashes is created by omitting a dot from the domain
CO    Character Omission Omitting a character from the domain
CS    Character Swap Swapping two consecutive characters in a domain
ACS   Adjacent Character Substitution replaces adjacent characters
ACI   Adjacent Character Insertion inserts adjacent character
HG    Homoglyphs replaces characters with characters that look similar
SP    Singular Pluralise creates a singular domain plural and vice versa
CR    Character Repeat Repeats a character of the domain name twice
DCR   Double Character Replacement repeats a character twice
CM    Common Misspellings are created from common misspellings
HP    Homophones Typos are created from sets of words that sound the same
VS    Vowel Swapping is created by swaps vowels
BF    Bitsquatting relies on random bit-errors to redirect connections
WTLD  Wrong Top Level Domain
WSLD  Wrong Second Level Domain
NS    Numeral Swap numbers, words and vice versa

## Extra Functions

MX      Checking for DNS's MX records
TXT     Checking for DNS's TXT records
IP      Checking for IP address
NS      Checks DNS NS records
CNAME   Checks DNS CNAME records
IDNA    Show international domain name


### TODO 

* GeoIp Lookup.
* Estimate popularity of a domain variant via google search
* Lookup whois record
* Checks for webpage similarity
* Distribute extra functions to workers







## Authors

* [Rangertaha](https://github.com/rangertaha)


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
