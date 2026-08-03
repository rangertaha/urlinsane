---
title: Bibliography
parent: Reference
nav_order: 5
---

# Bibliography
{: .no_toc }

- TOC
{:toc}

Typosquatting is measured better than most abuse categories, and the reason is
structural. The candidate space is *enumerable* — given a name and a set of
error models you can write down every variant — and the ground truth is
*queryable*, because a domain is either registered or it is not, and DNS,
WHOIS, certificate transparency logs and zone files will tell you which. A
researcher can therefore generate the whole neighbourhood of a name, look up
every member of it, and report a census rather than a sample. Almost every
algorithm in URLInsane traces back to a paper below that did exactly that and
published the error model it used.

{: .warning }
The PDFs in `docs/papers/` are byte-damaged: every non-ASCII byte in them was
replaced with the UTF-8 replacement character at some point before they were
committed, which breaks the compressed streams. `pdftotext` and Ghostscript both
report a broken cross-reference table on 25 of the 27 files, and only
`3663569.pdf` and `Measuring and Analyzing Typosquatting.pdf` still open. The
local links below are kept because the filenames are how these works are
referred to in the code and in `docs/REFS.md`, but use the original URL where
one is given.

## Which paper is behind which algorithm

The algorithm set is not a survey of what other tools do; each generator exists
because a measurement study found the technique in the wild. This table is the
mapping, and it is the answer to "why is this one here?".

| Algorithm | Paper | In `papers/` |
|---|---|---|
| `bf` | Dinaburg, *Bitsquatting* (Black Hat 2011) — bit errors in RAM reach DNS at internet scale | [BH_US_11_Dinaburg_Bitsquatting_WP.pdf](../../papers/BH_US_11_Dinaburg_Bitsquatting_WP.pdf) |
| `bf` | Schultz, *Examining the Bitsquatting Attack Surface* (DEF CON 21) | [DEFCON-21-Schultz…pdf](../../papers/DEFCON-21-Schultz-Examining-the-Bitsquatting-Attack-Surface-WP.pdf) |
| `cb` | Kintis et al., *Hiding in Plain Sight* (CCS 2017) — combosquatting outnumbers and outlives typo variants | [p569-kintisA.pdf](../../papers/p569-kintisA.pdf) |
| `co` `cs` `cr` `acs` `aci` | Szurdi et al., *The Long "Taile" of Typosquatting* (USENIX Sec 2014) | [sec14-paper-szurdi.pdf](../../papers/sec14-paper-szurdi.pdf) |
| `hs` | Nikiforakis et al., *Soundsquatting* — homophones as a squatting vector | see below |
| `xhs` | Valentim et al., *X-squatter* (ACM TOPS 2024) — sound-squatting **across** languages, ~15% of candidates carry TLS certificates | [3663569.pdf](../../papers/3663569.pdf) |
| `hr` | Unicode UTS #39, plus the IDN homograph literature | — |
| `cm` | Birkbeck / human spelling-error corpora | — |
| `afx` `nsc` `sep` | Duan et al. (NDSS 2021) and the PyPI/npm supply-chain measurements | [ndss2021_1B-1_23055_paper.pdf](../../papers/ndss2021_1B-1_23055_paper.pdf) |
| `tld` `sld` `tli` `fsd` | Agten et al. and the ccTLD/level measurement work; `fsd` narrows to the public suffix list's private section | [imc17-final215.pdf](../../papers/imc17-final215.pdf) |
| `tos` | ail-typo-squatting's ChangeOrder; word-order confusion in package names | — |

## Start here

Three works carry most of the weight. Between them they establish that the
candidate space can be enumerated exhaustively, that error sources include ones
no human generates, and that the highest-volume abuse involves no typing error
at all.

**The Long "Taile" of Typosquatting Domain Names.** Janos Szurdi, Balazs Kocso,
Gabor Cseh, Jonathan Spring, Mark Felegyhazi, Chris Kanich. USENIX Security
2014.
[PDF](../../papers/sec14-paper-szurdi.pdf) ·
[usenix.org](https://www.usenix.org/system/files/conference/usenixsecurity14/sec14-paper-szurdi.pdf)
A census of typo registrations across the whole of `.com` rather than just the
popular head. About half of the lexically-identified typo candidates are true
typo domains, and the paper estimates 20% of all `.com` registrations are typo
domains — most of them targeting the long tail, with only 6.8% aimed at the top
10,000 sites.

**Bitsquatting: DNS Hijacking without Exploitation.** Artem Dinaburg. Black Hat
USA 2011 white paper. [PDF](../../papers/BH_US_11_Dinaburg_Bitsquatting_WP.pdf)
Registers domains one *bit* away from frequently-resolved names and logs the
HTTP requests that arrive. Six months of logs show that random bit errors in
RAM — from manufacturing defects, heat and radiation — reach DNS often enough
to be exploitable, and that virtually every operating system and platform is
affected. This is why URLInsane's `bf` algorithm exists.

**Hiding in Plain Sight: A Longitudinal Study of Combosquatting Abuse.**
Panagiotis Kintis, Najmeh Miramirkhani, Charles Lever, Yizheng Chen, Roza
Romero-Gómez, Nikolaos Pitropakis, Nick Nikiforakis, Manos Antonakakis. ACM CCS
2017. [PDF](../../papers/p569-kintisA.pdf) ·
[acmccs.github.io](https://acmccs.github.io/papers/p569-kintisA.pdf)
Six years and 468 billion DNS records. Almost 60% of abusive combosquatting
domains live longer than 1,000 days, the volume grows year over year, and the
abuse spans phishing, social engineering, affiliate fraud, trademark abuse and
APT activity. No typing error is involved in any of it.

## Papers

### Measurement studies

**Seven Months' Worth of Mistakes: A Longitudinal Study of Typosquatting
Abuse.** Pieter Agten, Wouter Joosen, Frank Piessens, Nick Nikiforakis. NDSS
2015. [PDF](../../papers/01_3_1.pdf) ·
[ndss-symposium.org](https://www.ndss-symposium.org/wp-content/uploads/2017/09/01_3_1.pdf)
The first content-based longitudinal study: the typo domains of the top 500
sites, visited daily for seven months. 95% of popular domains are actively
targeted, few trademark owners register defensively, and squatted domains
change hands over time.

**Measuring and Analyzing Typosquatting Toward Fighting Abusive Domain
Registrations.** Janos Szurdi. PhD thesis, Carnegie Mellon University, July
2020. [PDF](../../papers/Measuring%20and%20Analyzing%20Typosquatting.pdf)
Collects the author's typosquatting work — including the USENIX and IMC papers
above and below — into one treatment covering measurement, economics, ethics
and intervention.

**Large-Scale Analysis of Pop-Up Scam on Typosquatting URLs.** Tobias Dam,
Lukas Daniel Klausner, Damjan Buhov, Sebastian Schrittwieser. ARES 2019.
[PDF](../../papers/1906.10762.pdf) ·
[arXiv:1906.10762](https://eprints.cs.univie.ac.at/7031/1/1906.10762.pdf)
Crawls typo domains derived from the Alexa top 1M and finds 9,857 pop-up
messages on 8,255 distinct URLs, 8,828 of them malicious. Most URLs served the
scam to one specific user agent only — a reminder that a single-fingerprint
crawl undercounts.

**A Smörgåsbord of Typos: Exploring International Keyboard Layout
Typosquatting.** Victor Le Pochat, Tom Van Goethem, Wouter Joosen. WTMC 2019
(IEEE Security and Privacy Workshops).
[PDF](../../papers/smorgasbord-wtmc19.pdf) ·
[lepoch.at](https://lepoch.at/files/smorgasbord-wtmc19.pdf)
Previous work assumed the US English layout. This examines typo domains on
non-US layouts for 100,000 popular domains, finds German users the most
targeted with over 15,000 registered typo domains, and finds defensive
registration patchy where it exists at all. This is the argument for
URLInsane's per-layout keyboard model in `pkg/kb`.

**Poster: A Smörgåsbord of Typos.** Victor Le Pochat, Tom Van Goethem, Wouter
Joosen. IEEE S&P 2019 poster session.
[PDF](../../papers/hotcrp_sp19posters-final27.pdf) ·
[ieee-security.org](https://www.ieee-security.org/TC/SP2019/posters/hotcrp_sp19posters-final27.pdf)
The two-page version of the above.

**A User Study of the Effectiveness of Typosquatting Techniques.** Jeffrey
Spaulding, Ah Reum Kang, Shambhu Upadhyaya, Aziz Mohaisen. IEEE CNS 2016
(poster). [PDF](../../papers/CNS_typosquatting.pdf) · original URL in
`docs/REFS.md`, now dead

**Understanding the Effectiveness of Typosquatting Techniques.** Jeffrey
Spaulding, DaeHun Nyang, Aziz Mohaisen. HotWeb 2018.
[PDF](../../papers/hotweb18a.pdf) ·
[cs.ucf.edu](http://www.cs.ucf.edu/~mohaisen/doc/hotweb18a.pdf)

### Specific attack classes

**Email Typosquatting.** Janos Szurdi, Nicolas Christin. ACM IMC 2017.
[PDF](../../papers/imc17-final215.pdf) ·
[sigcomm.org](https://conferences.sigcomm.org/imc/2017/papers/imc17-final215.pdf)
Registers 76 typo domains and collects the mail sent to them for seven months.
Extrapolates roughly 3,585 misdirected legitimate emails per year across three
domains, some containing visa documents and medical records. This is the
evidence behind treating an MX record on a squat as a materially worse finding
than a parked page.

**Examining the Bitsquatting Attack Surface.** Jaeson Schultz, Cisco. DEF CON
21 white paper, 2013.
[PDF](../../papers/DEFCON-21-Schultz-Examining-the-Bitsquatting-Attack-Surface-WP.pdf) ·
[defcon.org](https://defcon.org/images/defcon-21/dc-21-presentations/Schultz/DEFCON-21-Schultz-Examining-the-Bitsquatting-Attack-Surface-WP.pdf)
Extends Dinaburg: describes previously unknown forms of bitsquatting — notably
bit flips in URL *delimiters*, which reach domains that are otherwise
unsquattable — and proposes mitigations that do not require mass defensive
registration.

**X-squatter: AI Multilingual Generation of Cross-Language Sound-squatting.**
Rodolfo Vieira Valentim, Idilio Drago, Marco Mellia, Federico Cerutti. ACM
Transactions on Privacy and Security 27(3), Article 21, June 2024.
[PDF](../../papers/3663569.pdf) ·
[doi:10.1145/3663569](https://doi.org/10.1145/3663569)
Generates sound-squatting candidates across languages with a transformer, then
checks them against hundreds of millions of TLS certificates: roughly 15% of
generated sound-squats have certificates, against 7% for other squatting types.
Also finds hundreds of sound-squat candidates in three years of PyPI package
history. The cross-language framing is directly relevant to URLInsane's `hs`
algorithm and its multilingual datasets.

**Towards Measuring Supply Chain Attacks on Package Managers for Interpreted
Languages.** Ruian Duan, Omar Alrawi, Ranjita Pai Kasturi, Ryan Elder, Brendan
Saltaformaggio, Wenke Lee. NDSS 2021.
[PDF](../../papers/ndss2021_1B-1_23055_paper.pdf)
Compares the security features of package managers for interpreted languages
and applies metadata, static and dynamic analysis to registry abuse, reporting
339 new malicious packages. The reference point for the package-registry side
of the [named-entity surface]({{ site.baseurl }}/attack/surface/).

**The Wolf of Name Street: Hijacking Domains Through Their Nameservers.**
Thomas Vissers, Timothy Barron, Wouter Joosen, Nick Nikiforakis. ACM CCS 2017.
[PDF](../../papers/p957-vissersA.pdf) ·
[acmccs.github.io](https://acmccs.github.io/papers/p957-vissersA.pdf)
Adjacent rather than squatting proper: takeover through the nameserver rather
than through the name. Included because a scan that resolves NS records is
looking at the same infrastructure.

**Typosquat Cyber Crime Attack Detection via Smartphone.** Z. Zulkefli, M. M.
Singh, A. R. Mohd Shariff, A. Samsudin. Procedia Computer Science 124 (2017)
664–671.
[PDF](../../papers/Typosquat_Cyber_Crime_Attack_Detection_via_Smartph.pdf) ·
[doi:10.1016/j.procs.2017.12.203](https://doi.org/10.1016/j.procs.2017.12.203)

### Detection and modelling

**It's All in the Name: Why Some URLs are More Vulnerable to Typosquatting.**
Rashid Tahir, Ali Raza, Faizan Ahmad, Jehangir Kazi, Fareed Zaffar, Chris
Kanich, Matthew Caesar. IEEE INFOCOM 2018.
[PDF](../../papers/tahir2018itsall.pdf) ·
[cs.uic.edu](https://www.cs.uic.edu/~ckanich/papers/tahir2018itsall.pdf)
Models the relationship between hand anatomy, keyboard layout and typing error
to compute a per-URL "Hardness-Quotient" — a likelihood of being mistyped — and
predicts the most likely typos for defensive registration. The closest thing in
the literature to a principled ranking function over generated candidates.

**Harvesting SSL Certificate Data to Identify Web-Fraud.** Mishari Almishari,
Emiliano De Cristofaro, Karim El Defrawy, Gene Tsudik. International Journal of
Network Security 14(6), 324–338, November 2012.
[PDF](../../papers/IJSN12.pdf) ·
[emilianodc.com](https://emilianodc.com/PAPERS/IJSN12.pdf)
Builds a classifier for fraudulent domains — phishing and typosquatting — from
the properties of their SSL certificates. An early version of the argument that
certificate data is a usable signal over squatted names.

**Deepsquatting: Learning-Based Typosquatting Detection at Deeper Domain
Levels.** Paolo Piredda, Davide Ariu, Battista Biggio, Igino Corona, Luca
Piras, Giorgio Giacinto, Fabio Roli. AI*IA 2017, LNCS 10640.
[PDF](../../papers/piredda17-AIIA.pdf)
Learns a similarity measure between domain names to detect typosquatting in DNS
traffic, including at levels below the registrable domain.

**DNS Typo-squatting Domain Detection: A Data Analytics & Machine Learning
Based Approach.** Abdallah Moubayed, MohammadNoor Injadat, Abdallah Shami,
Hanan Lutfiyya. arXiv:2012.13604, December 2020.
[PDF](../../papers/2012.13604v1.pdf) ·
[arXiv](https://arxiv.org/pdf/2012.13604)
Eight domain-name-derived features, a majority-voting ensemble of five
classifiers, and K-means clustering used to validate the same features on
unlabelled data.

### Linguistic and NLP resources

**GitHub Typo Corpus: A Large-Scale Multilingual Dataset of Misspellings and
Grammatical Errors.** Masato Hagiwara, Masato Mita. LREC 2020, 6761–6768.
[PDF](../../papers/2020.lrec-1.835.pdf) ·
[aclanthology.org](https://aclanthology.org/2020.lrec-1.835.pdf)
More than 350k edits and 65M characters across 15+ languages, harvested from
git history. The kind of observed-error corpus that a misspelling algorithm
should be trained on rather than hand-written — see [linguistic
datasets]({{ site.baseurl }}/guide/datasets/).

### Legal, policy and industry reports

**Cybersquatting, Typosquatting, and Domaining: Ten Years Under the
Anti-Cybersquatting Consumer Protection Act.** Carl C. Butzer, Jason P.
Reinsch. Law review article, 2009. [PDF](../../papers/1276.pdf) ·
[jw.com](https://www.jw.com/wp-content/uploads/2016/05/1276.pdf)
The legal remedy side: what the ACPA does and does not reach.

**A Study of Whois Privacy and Proxy Service Abuse.** NPL Management Ltd for
ICANN, 20 September 2013; primary author Richard Clayton, University of
Cambridge. [PDF](../../papers/pp-abuse-study-20sep13-en.pdf) ·
[gnso.icann.org](https://gnso.icann.org/sites/default/files/filefield_41831/pp-abuse-study-20sep13-en.pdf)
Measures how much more often domains used for illegal or harmful activity sit
behind privacy or proxy registration than domains generally. Relevant to
reading WHOIS output on a squat.

**Typosquatting – A New Menace to Society.** Palak Sharma. International
Journal of Creative Research Thoughts 10(5), May 2022.
[PDF](../../papers/IJCRT2205845.pdf) ·
[ijcrt.org](https://ijcrt.org/papers/IJCRT2205845.pdf)
A survey of harms to victims and of the gap between those harms and the legal
provisions available, with an Indian focus.

**An investigation of phishing awareness and education over time: When and how
to best remind users.** Benjamin Reinheimer, Lukas Aldag, Peter Mayer, Mattia
Mossano, Reyhan Duezguen, Bettina Lofthouse, Tatiana von Landesberger, Melanie
Volkamer. SOUPS 2020. [PDF](../../papers/soups2020-reinheimer_0.pdf) ·
[usenix.org](https://www.usenix.org/system/files/soups2020-reinheimer_0.pdf)
Not about squatting, but about the other half of the problem: how long user
training against lookalike names actually lasts.

**2024 Data Breach Investigations Report.** Verizon, 2024.
[PDF](../../papers/2024-dbir-data-breach-investigations-report.pdf)
Industry breach statistics. No URL is recorded for it in `docs/REFS.md`.

**`Final-Paper-cyse494-copy.pdf`** — the title could not be extracted. The
document metadata gives the author as "Chibuike, Oga" and a creation date of
23 June 2023, and the bibliography links point at phishing and information-
security literature, but the file is damaged and no title is recoverable from
it. There is no URL for it in `docs/REFS.md`.

### Referenced in REFS.md without a local copy

`docs/REFS.md` also lists these, for which no PDF was archived here:

- **Defending Against Typosquatting Attacks In Programming Language-Based
  Package Repositories.** Matthew Taylor. MS thesis, University of Kansas, May
  2020.
  [kuscholarworks.ku.edu](https://kuscholarworks.ku.edu/server/api/core/bitstreams/7bddaa7e-59b4-40af-a73b-5d7c93d28803/content)
- **TypoAlert: a browser extension against typosquatting.** Francesco Blefari,
  Angelo Furfaro, Giovambattista Ianni, Alessandro Viscomi. SEBD 2024.
  [sebd2024.unica.it](https://sebd2024.unica.it/papers/paper85.pdf)
- An ARES 2016 paper at
  [cs.ucf.edu/~mohaisen/doc/ares16.pdf](https://www.cs.ucf.edu/~mohaisen/doc/ares16.pdf)
  — the link no longer resolves and no local copy exists, so it is unidentified.
- **Typosquatting Domains Analysis.** Recorded Future blog.
  [recordedfuture.com](https://www.recordedfuture.com/blog/typosquatting-domains-analysis)
- **Python Typosquatting for Fun not Profit.** William Bengtson.
  [medium.com](https://medium.com/@williambengtson/python-typosquatting-for-fun-not-profit-99869579c35d)

## Related tools

Prior art, and the tools URLInsane is usually compared against.

| Tool | What it does |
|---|---|
| [urlcrazy](https://github.com/urbanadventurer/urlcrazy) | The original Ruby generator: typo variants of a domain plus DNS and popularity checks. URLInsane's README names it as one of the tools this project was inspired by. |
| [dnstwist](https://github.com/elceef/dnstwist) | The most widely used Python tool. Generates permutations and resolves them, with WHOIS, GeoIP, banner grabbing, MX detection and fuzzy page comparison. |
| [DomainFuzz](https://github.com/monkeym4ster/DomainFuzz) | Node.js domain-name permutation and registration checking. |
| [ail-typo-squatting](https://github.com/typosquatter/ail-typo-squatting) | Python library that generates typo variants from a domain, usable as a library rather than only a CLI. |
| [pypi-squatting](https://github.com/typosquatter/pypi-squatting) | The same idea aimed at PyPI package names rather than domains. |

## Data sources

The datasets behind the algorithms, as opposed to the literature behind them.

- **[kbdlayout.info](http://kbdlayout.info/)** — keyboard layout definitions,
  including key geometry. The source for the layouts in
  [`pkg/kb`]({{ site.baseurl }}/KB/), which is what makes adjacency
  layout-specific rather than QWERTY-only.
- **[Birkbeck spelling error corpora](https://www.dcs.bbk.ac.uk/~roger/corpora.html)**
  — Roger Mitton's collections of *observed* human spelling errors. Observed
  error data, not generated: the difference matters, because a generator tuned
  on real errors ranks differently from one tuned on edit distance.
- **[MaxMind GeoIP](https://www.maxmind.com/en/geoip-databases)** — IP to
  country/ASN mapping, used to say where a resolved squat is hosted.
- **[Public Suffix List](https://publicsuffix.org/)** — the authoritative list
  of registrable-domain boundaries. Required to tell a registrable name from a
  subdomain, which is what separates a squat from a `levelsquat`.

---

Back to **[Reference](../)**.
