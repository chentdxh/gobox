# gobox 🧰

All-in-one terminal toolbox for developers. One binary, zero dependencies, pipe-friendly.

## Commands

| Command | Alias | What it does |
|---------|-------|-------------|
| `json` | `jq` | Format & indent JSON |
| `jsonc` | — | Compact JSON to one line |
| `b64e` | `b64enc` | Base64 encode |
| `b64d` | `b64dec` | Base64 decode |
| `urle` | `urlenc` | URL encode |
| `urld` | `urldec` | URL decode |
| `hash` | — | Hash (md5, sha1, sha256, sha512) |
| `ts` | `timestamp` | Unix timestamp → date |
| `d2t` | `date2ts` | Date string → unix timestamp |
| `hexe` | `hexenc` | Hex encode |
| `hexd` | `hexdec` | Hex decode |
| `jwt` | — | Decode JWT header + payload |
| `count` | `wc` | Count chars, words, lines |
| `all` | — | Show all transforms at once |

## Install

```bash
go install github.com/chentdxh/gobox/cmd/gobox@latest
```

## Usage

Every command accepts **piped input** or **direct arguments**:

```bash
# Pipe JSON from curl
curl -s https://api.example.com/user | gobox json

# Base64 encode
gobox b64e "hello world"

# Hash a string
gobox hash sha256 "password123"

# Decode JWT
gobox jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0In0.xxx

# Convert timestamp to date
gobox ts 1700000000

# See everything at once
echo "hello" | gobox all
```
