# URL Insane Api Requirements

OPTIONS /
Response:

```json
{
    "domain": {
        "type": "input",
        "description": "The domain",
        "optional": false
    },
    "funcs": {
        "type": "multi-select",
        "description": "Extra functions for data or filtering (default [idna])",
        "optional": true,
        "values": [
            //
        ],
    },
    "keyboards": {
        "type": "multi-select",
        "description": "Keyboards/layouts ID to use (default [en1])",
        "optional": true,
        "values": [
            {
                "value": "AR2",
                "name": "Arabic...",
                "description": "Arabic PC keyboard layout"
            }
        ]
    },
    "typos": {
        "type": "multi-select",
        "description": "Types of typos to perform (default [all])",
        "optional": true,
        "values": [
            //

        ]
    }
}
```

POST /
body:

```json
{
    "domains": ["google.com"],
    "funcs": [
        "all"
    ],
    "keyboards": [
        "en1"
    ],
    "typos": [
        "co"
    ],
}
```

response:

```json
{
    "headers": [
        "Suffix",
        "IPv4",
        "CNAME",
        "GEO",
        "IDNA",
        "Typo",
        "Variant",
        "Type",
        "IPv6",
        "Redirect",
        "TXT",
        "NS",
        "SIM",
        "Live",
        "MX"
    ],
    "rows": [
        {
            "CNAME": "",
            "GEO": "United States",
            "IDNA": "oogle.com",
            "IPv4": "104.28.28.162\n104.28.29.162",
            "IPv6": "2606:4700:30::681c:1ca2\n2606:4700:30::681c:1da2",
            "Live": true,
            "MX": "mx.zoho.com\nmx2.zoho.com",
            "NS": "amir.ns.cloudflare.com\ngwen.ns.cloudflare.com",
            "Redirect": "oogle.com",
            "SIM": "",
            "Suffix": "",
            "TXT": "v=spf1 +a +mx +ip4:204.9.184.9 +include:zoho.com ~all",
            "Type": "",
            "Typo": {
                "code": "CO",
                "name": "Character Omission",
                "description": "Character Omission Omitting a character from the domain."
            },
            "Variant": "oogle.com"
        },
        {
            "CNAME": "",
            "GEO": "United States",
            "IDNA": "googl.com",
            "IPv4": "172.217.7.4",
            "IPv6": "2607:f8b0:4006:81b::2004",
            "Live": true,
            "MX": "",
            "NS": "ns1.google.com\nns3.google.com\nns2.google.com\nns4.google.com",
            "Redirect": "www.google.com",
            "SIM": "72%",
            "Suffix": "",
            "TXT": "v=spf1 -all",
            "Type": "",
            "Typo": {
                "code": "CO",
                "name": "Character Omission",
                "description": "Character Omission Omitting a character from the domain."
            },
            "Variant": "googl.com"
        },

    ]
}
```
