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



## Installation

To get the latest updates, create the binary executable with the **make** command or 
[download](https://github.com/cybersectech-org/urlinsane/releases/tag/0.3.0) one of the 
pre-built release binaries. 

Get the project
```bash
go get github.com/cybersectech-org/urlinsane
```

Go to the project folder and run the **make** command.
```bash
cd ~/go/src/github.com/cybersectech-org/urlinsane/
make
```

After building the binary you can execute it within the **builds** directory 
that was created by the **make** command. 
```bash
cd builds/
./urlinsane -h
```


## Usage
Generate variations for `google.com` using the character omission **(CO)** algorithm.
```
urlinsane typo google.com -t co
 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: 0.3.0

  LIVE | TYPE |   TYPO    | SUFFIX |   IDNA     
+------+------+-----------+--------+-----------+
       | CO   | oogle.com | com    | oogle.com  
       | CO   | gogle.com | com    | gogle.com  
       | CO   | goole.com | com    | goole.com  
       | CO   | googl.com | com    | googl.com  
       | CO   | gogle.com | com    | gogle.com  
       | CO   | googe.com | com    | googe.com 
```


Generate variations for `google.com` using the character omission **(CO)** algorithm and check for **ip** addresses. 
```
urlinsane typo google.com -t co -x ip

 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: 0.3.0

   LIVE  | TYPE |   TYPO    | SUFFIX |      IPV4      |           IPV6            
+--------+------+-----------+--------+----------------+--------------------------+
  ONLINE | CO   | oogle.com | com    | 104.28.28.162  | 2606:4700:30::681c:1ca2   
         |      |           |        | 104.28.29.162  | 2606:4700:30::681c:1da2   
  ONLINE | CO   | gogle.com | com    | 172.217.10.68  | 2607:f8b0:4004:80a::2004  
  ONLINE | CO   | gogle.com | com    | 172.217.10.68  | 2607:f8b0:4004:80a::2004  
  ONLINE | CO   | googl.com | com    | 172.217.10.132 | 2607:f8b0:4004:800::2004  
  ONLINE | CO   | goole.com | com    | 217.160.0.201  |                           
  ONLINE | CO   | googe.com | com    | 162.243.10.151 |   
```

Generate variations for `google.com` using the character omission **(CO)** algorithm. 
 Also execute extra functions to get the **ip** addresses, **idna** format and check for **ns** records. 
```
urlinsane typo google.com -t co -x ip -x idna -x ns

 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: 0.3.0

   LIVE  | TYPE |   TYPO    | SUFFIX |      IPV4      |           IPV6           |   IDNA    |        NS         
+--------+------+-----------+--------+----------------+--------------------------+-----------+------------------+
  ONLINE | CO   | gogle.com | com    | 172.217.10.68  | 2607:f8b0:4004:80a::2004 | gogle.com |                   
  ONLINE | CO   | oogle.com | com    | 104.28.29.162  | 2606:4700:30::681c:1da2  | oogle.com | mx.zoho.com       
         |      |           |        | 104.28.28.162  | 2606:4700:30::681c:1ca2  |           | mx2.zoho.com      
  ONLINE | CO   | gogle.com | com    | 172.217.10.68  | 2607:f8b0:4004:80a::2004 | gogle.com |                   
  ONLINE | CO   | googl.com | com    | 172.217.10.132 | 2607:f8b0:4004:800::2004 | googl.com |                   
  ONLINE | CO   | goole.com | com    | 217.160.0.201  |                          | goole.com | mx01.1and1.co.uk  
         |      |           |        |                |                          |           | mx00.1and1.co.uk  
  ONLINE | CO   | googe.com | com    | 162.243.10.151 |                          | googe.com |                   

```

For more details look at the **-h --help** output.
```
urlinsane typo -h

    
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

* Open Source
* Binary executable, written in Go with no dependencies. 
* Will have all the functionality of URLCrazy and DNSTwist. 
* Contains 19 typosquatting algorithms and 7 extra functions to retrieve additional data
* Modular architecture for language, keyboard, typo algorithm, and functions extensibility.
* Supports multiple keyboard layouts found in English, Spanish, Russian, Finnish, and Arabic.
* Supports multiple languages with the ability to add more languages with ease.
* Concurrent function (**-x --funcs**) workers to retrieve additional info on each record.
* Concurrent typo squatting workers.


## Keyboards

  AR3	Arabic North African keyboard layout
  EN1	English QWERTY keyboard layout
  EN4	English DVORAK keyboard layout
  AR1	Arabic keyboard layout
  EN2	English AZERTY keyboard layout
  EN3	English QWERTZ keyboard layout
  FI1	Finnish QWERTY keyboard layout
  RU2	Phonetic Russian keyboard layout
  RU3	PC Russian keyboard layout
  ES1	Spanish keyboard layout
  ES2	Spanish ISO keyboard layout
  AR2	Arabic PC keyboard layout
  AR4	Arabic keyboard layout
  RU1	Russian keyboard layout
  ALL	Use all keyboards
  
## Languages

URLInsane supports multiple languages with the ability to add more languages 
with ease. If you know where I can find a list of commonly misspelled words or 
homophones for languages other than English please send the info to 
rangertaha@gmail.com or add it yourself and submit a pull request. 


### English

* Over 8000 common misspellings
* Over 500 common homophones
* English alphabet, vowels, homoglyphs, and numerals
* Common keyboard layouts (qwerty, azerty, qwertz, dvorak)

### Spanish

* Spanish alphabet, vowels, homoglyphs, and numerals
* Basic and ISO keyboard layouts 

### Russian

* Russian alphabet, vowels, homoglyphs, and numerals
* Common 3 keyboard layouts

### Finish

* Finish alphabet, vowels, homoglyphs, and numerals
* Finnish QWERTY keyboard layout

### Arabic

* Arabic alphabet, homoglyphs, and numerals
* 4 Common keyboard layouts

### Persian

* Persian alphabet and numerals
* One Common keyboard layout


### Hebrew

* Hebrew alphabet
* One Common keyboard layout

### Armenian

* Armenian alphabet and numerals
* 2 Common keyboard layouts


## Algorithms

The modular architecture for code extensibility allows developers to add new typosquatting
algorithms with ease. Currently, we have implemented 19 
typosquatting algorithms. 

1. **Missing Dot(MD)** is created by omitting a dot from the domain.
2. **Missing Dashes(MDS)** is created by omitting a dash from the domain.
3. **Strip Dashes(SD)** is created by stripping all dashes from the domain
4. **Character Omission(CO)** Omitting a character from the domain
5. **Character Swap(CS)** swaps two consecutive characters in a domain
6. **Adjacent Character Substitution(ACS)** replaces adjacent characters
7. **Adjacent Character Insertion(ACI)** inserts adjacent character
8. **Homoglyphs(HG)** replaces characters with characters that look similar
9. **Singular Pluralise(SP)** creates a singular domain plural and vice versa
10. **Character Repeat(CR)** repeats a character of the domain name twice
11. **Double Character Replacement(DCR)** repeats a character twice
12. **Common Misspellings(CM)** are created from common misspellings
13. **Homophones(HP)** Typos are created from sets of words that sound the same
14. **Vowel Swapping(VS)** is created by swaps vowels
15. **Bitsquatting(BF)** relies on random bit-errors
16. **Wrong Top Level Domain(WTLD)**
17. **Wrong Second Level Domain(W2TLD)**
18. **Numeral Swap(NS)** numbers, words and vice versa
19. **Wrong Third Level Domain(W3TLD)**

## Extra Functions

- **IDNA**    Show international domain name (Default)
- **MX**      Checking for DNS's MX records
- **TXT**     Checking for DNS's TXT records
- **IP**      Checking for IP address
- **NS**      Checks DNS NS records
- **CNAME**   Checks DNS CNAME records
- **SIM**     Show domain similarity % using fuzzy hashing with ssdeep
- **301**	Show domains redirects
- **GEO**	Show domains redirects

## TODO 

* Complete tool comparison
* Extract keywords from domains. Keywords will be used for additional algorithms
* Estimate popularity of a domain variant via google search
* Lookup whois record
* Emoji domains
* Grabs HTTP and SMTP service banners



## Tool Comparisons


|      **Algorithms**             | URLInsane | URLCrazy  | DNSTwist   | DomainFuzz  |
|                                 |           |           |            |             |          
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
| Cardinal Number Swap            |           |           |            |             |           
| Hyphenation                     |           |           |      X     |             |         
| Combosquatting(Keywords)        |           |           |            |             |           
| Multithreaded Algorithms        |     X     |     ?     |      X     |      X      |         



|      **Extra Functions**            | URLInsane  | URLCrazy  | DNSTwist   | DomainFuzz | 
|-------------------------------------|-----------|-----------|------------|-------------|
|                                     |           |           |            |             |           
| Live/Online Check                   |     X     |     X     |      X     |             |           
| Google Popularity Estimate          |           |     X     |            |             |           
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
| HTTP/SMTP Banner                    |           |           |      X     |             |           
| WHOIS Info                          |           |           |      X     |             |           
| Test MX email intercepts            |           |           |      X     |             |           
| Multithreaded Extra Functions       |     X     |           |      X     |      X      |           


| **Language (Keyboards)** | URLInsane  | URLCrazy  | DNSTwist   | DomainFuzz  | XN-Twist |
|--------------------------|-----------|-----------|------------|-------------|-----------|
|                          |           |           |            |             |           |
| Arabic (4)               |     X     |           |            |             |           |
| Armenian (3)             |     X     |           |            |             |           |
| English (4)              |     X     |     X     |      X     |      X      |     X     |
| Finnish (1)              |     X     |           |            |             |           |
| Russian (3)              |     X     |           |            |             |           |
| Spanish (2)              |     X     |           |            |             |           |
| Hebrew (1)               |     X     |           |            |             |           |
| Persian (1)              |     X     |           |            |             |           |


## Authors

* [Rangertaha](https://github.com/rangertaha)


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
