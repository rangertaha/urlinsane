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
        "live",
        "type",
        "typo",
        "suffix",
        "ipv4",
        "ipv6",
        "idna",
        "ns"
    ],
    "rows": [
        {
            "live": "ONLINE",
            "type": "",
            "typo": "",
            "suffix": "",
            "ipv4": "",
            "ipv6": "",
            "idna": "",
            "ns": ""
        }
    ]
}
```
