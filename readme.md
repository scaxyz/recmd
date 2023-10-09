# Record/Replay commands

A simple lib (with cli) to record the outputs of an command with the same timedelays.

## Usage lib
### record
```go
recorder := recmd.NewRecorder()
record := recorder.Record("curl duckduckgo.com",nil)
record2 := recorder.Record("bash", os.Stdin)

// records are stored as a json

```

### read
```go
var record recmd.Record
json.Unmarschal(data, &record)

reader := record.Reader();
data, _ := io.ReadAll(reader);
```


## Usage cli

### record 

```bash
recmd record wget duckduckgo.com

# or

recmd record --stdin bash

# or 

recmd record --stdin -o apt-record.json apt update

# or

recmd record --plain-text bash

```

### replay

```bash

recmd replay some-record.json

# or 

recmd replay --no-delays some-record.json

```

### convert

```bash

recmd conv-plain input.json output.json

# or

recmd ctp input.json
# => produces input.plain.json

```

## Examples cli
### `recmd record wget duckduckgo.com`
Produces `recmd-20230710_171624.json` with:
<details>

<summary>Click to expand</summary>

```json
{
    "command": "/usr/bin/wget duckduckgo.com",
    "out": {},
    "in": {},
    "err": {
        "160578848": "MzAxIE1vdmVkIFBlcm1hbmVudGx5CkxvY2F0aW9uOiBodHRwczovL2R1Y2tkdWNrZ28uY29tLyBbZm9sbG93aW5nXQo=",
        "163537869": "LS0yMDIzLTA3LTEwIDE3OjE2OjI1LS0gIGh0dHBzOi8vZHVja2R1Y2tnby5jb20vCg==",
        "202423050": "Q29ubmVjdGluZyB0byBkdWNrZHVja2dvLmNvbSAoZHVja2R1Y2tnby5jb20pfDQwLjExNC4xNzcuMTU2fDo0NDMuLi4g",
        "274610762": "Y29ubmVjdGVkLgo=",
        "454187285": "SFRUUCByZXF1ZXN0IHNlbnQsIGF3YWl0aW5nIHJlc3BvbnNlLi4uIA==",
        "5276292": "LS0yMDIzLTA3LTEwIDE3OjE2OjI0LS0gIGh0dHA6Ly9kdWNrZHVja2dvLmNvbS8K",
        "5402847": "UmVzb2x2aW5nIGR1Y2tkdWNrZ28uY29tIChkdWNrZHVja2dvLmNvbSkuLi4g",
        "546082587": "MjAwIE9LCg==",
        "546471465": "TGVuZ3RoOiA=",
        "546591196": "NjQ2OQ==",
        "546744357": "ICg2LDNLKQ==",
        "546866055": "IFt0ZXh0L2h0bWxdCg==",
        "549915916": "U2F2aW5nIHRvOiDigJhpbmRleC5odG1sLjHigJkK",
        "551886110": "CiAgICAgMEs=",
        "552103212": "IA==",
        "552212133": "Lg==",
        "552315672": "Lg==",
        "552418116": "Lg==",
        "552524321": "Lg==",
        "552660157": "Lg==",
        "552766094": "Lg==",
        "552884054": "IA==",
        "552997876": "IA==",
        "553107625": "IA==",
        "553287710": "IA==",
        "553387922": "IA==",
        "553498007": "IA==",
        "553607066": "IA==",
        "553714034": "IA==",
        "553822701": "IA==",
        "553928863": "IA==",
        "554038708": "IA==",
        "554145595": "IA==",
        "554251486": "IA==",
        "554357937": "IA==",
        "554465780": "IA==",
        "554572662": "IA==",
        "554679018": "IA==",
        "554789026": "IA==",
        "554895839": "IA==",
        "555005999": "IA==",
        "555125785": "IA==",
        "555282025": "IA==",
        "555385833": "IA==",
        "555497523": "IA==",
        "555606971": "IA==",
        "555715936": "IA==",
        "555822520": "IA==",
        "555933861": "IA==",
        "556045494": "IA==",
        "556252187": "IA==",
        "556354670": "IA==",
        "556472666": "IA==",
        "556588400": "IA==",
        "556700585": "IA==",
        "556809891": "IA==",
        "556921505": "IA==",
        "557107291": "ICAgICA=",
        "557278646": "ICAgIA==",
        "557449182": "ICAg",
        "557621808": "MTAwJQ==",
        "557803178": "IDksNTNN",
        "557976417": "PTAsMDAxcw==",
        "560675453": "Cgo=",
        "565019232": "MjAyMy0wNy0xMCAxNzoxNjoyNSAoOSw1MyBNQi9zKSAtIOKAmGluZGV4Lmh0bWwuMeKAmSBzYXZlZCBbNjQ2OS82NDY5XQoK",
        "6782279": "NDAuMTE0LjE3Ny4xNTYKQ29ubmVjdGluZyB0byBkdWNrZHVja2dvLmNvbSAoZHVja2R1Y2tnby5jb20pfDQwLjExNC4xNzcuMTU2fDo4MC4uLiA=",
        "88060736": "Y29ubmVjdGVkLgo=",
        "88089573": "SFRUUCByZXF1ZXN0IHNlbnQsIGF3YWl0aW5nIHJlc3BvbnNlLi4uIA=="
    }
}
```

</details>
<br>

### `recmd replay recmd-20230710_171624`
Displayes:
```text
Replaying:  /usr/bin/wget duckduckgo.com
--2023-07-10 17:16:24--  http://duckduckgo.com/
Resolving duckduckgo.com (duckduckgo.com)... 40.114.177.156
Connecting to duckduckgo.com (duckduckgo.com)|40.114.177.156|:80... connected.
HTTP request sent, awaiting response... 301 Moved Permanently
Location: https://duckduckgo.com/ [following]
--2023-07-10 17:16:25--  https://duckduckgo.com/
Connecting to duckduckgo.com (duckduckgo.com)|40.114.177.156|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 6469 (6,3K) [text/html]
Saving to: ‘index.html.1’

     0K ......                                                100% 9,53M=0,001s

2023-07-10 17:16:25 (9,53 MB/s) - ‘index.html.1’ saved [6469/6469]
```

### `recmd conv-plain recmd-20230710_171624.json`
Produces `recmd-20230710_171624.plain.json` with:

<details>
<summary>Click to expand</summary>


```json
{
    "command": "/usr/bin/wget duckduckgo.com",
    "err": {
        "160578848": "301 Moved Permanently\nLocation: https://duckduckgo.com/ [following]\n",
        "163537869": "--2023-07-10 17:16:25--  https://duckduckgo.com/\n",
        "202423050": "Connecting to duckduckgo.com (duckduckgo.com)|40.114.177.156|:443... ",
        "274610762": "connected.\n",
        "454187285": "HTTP request sent, awaiting response... ",
        "5276292": "--2023-07-10 17:16:24--  http://duckduckgo.com/\n",
        "5402847": "Resolving duckduckgo.com (duckduckgo.com)... ",
        "546082587": "200 OK\n",
        "546471465": "Length: ",
        "546591196": "6469",
        "546744357": " (6,3K)",
        "546866055": " [text/html]\n",
        "549915916": "Saving to: ‘index.html.1’\n",
        "551886110": "\n     0K",
        "552103212": " ",
        "552212133": ".",
        "552315672": ".",
        "552418116": ".",
        "552524321": ".",
        "552660157": ".",
        "552766094": ".",
        "552884054": " ",
        "552997876": " ",
        "553107625": " ",
        "553287710": " ",
        "553387922": " ",
        "553498007": " ",
        "553607066": " ",
        "553714034": " ",
        "553822701": " ",
        "553928863": " ",
        "554038708": " ",
        "554145595": " ",
        "554251486": " ",
        "554357937": " ",
        "554465780": " ",
        "554572662": " ",
        "554679018": " ",
        "554789026": " ",
        "554895839": " ",
        "555005999": " ",
        "555125785": " ",
        "555282025": " ",
        "555385833": " ",
        "555497523": " ",
        "555606971": " ",
        "555715936": " ",
        "555822520": " ",
        "555933861": " ",
        "556045494": " ",
        "556252187": " ",
        "556354670": " ",
        "556472666": " ",
        "556588400": " ",
        "556700585": " ",
        "556809891": " ",
        "556921505": " ",
        "557107291": "     ",
        "557278646": "    ",
        "557449182": "   ",
        "557621808": "100%",
        "557803178": " 9,53M",
        "557976417": "=0,001s",
        "560675453": "\n\n",
        "565019232": "2023-07-10 17:16:25 (9,53 MB/s) - ‘index.html.1’ saved [6469/6469]\n\n",
        "6782279": "40.114.177.156\nConnecting to duckduckgo.com (duckduckgo.com)|40.114.177.156|:80... ",
        "88060736": "connected.\n",
        "88089573": "HTTP request sent, awaiting response... "
    },
    "in": {},
    "out": {}
}
```

</details>
<br>

## Installation

```
go install github.com/scaxyz/recmd
```

## Known bugs
- when recording the `bash` executeable, typing `exit` and pressing `enter` requires a second `enter` to exit

## Some records are slower than the original command