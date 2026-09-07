# utility-belt

A fast, modular, offline, and stateless CLI tool designed for developers' everyday tasks.

## Features

- 🔐 **Password Generator:** Cryptographically secure passwords with configurable length and character classes (`crypto/rand`).
- 🆔 **UUID Generator:** Canonical UUID generation supporting UUID v4 (random) and UUID v7 (timestamp-based per RFC 9562).
- 🔢 **Base64 Encoder/Decoder:** Encode and decode data in standard RFC 4648 Base64 format with strict validation.
- 🔑 **JWT Decoder:** Inspect Header and Payload of any JSON Web Token (JWS compact format) as readable JSON, with expiration detection.
- ⚡ **Piping & Automation:** Clean `stdout` output for Unix pipelines (`pbcopy`, redirects, etc.) and diagnostic messages routed to `stderr`.

---

## Installation

### With `go install` (Go 1.21+)

```bash
go install github.com/kaioa1meida/utility-belt@latest.
```

### From Source

```bash
git clone https://github.com/kaioa1meida/utility-belt.git
cd utility-belt
go build -o utility-belt .
```

---

## Usage

### Password Generator

```bash
# Generate a standard 16-character password
utility-belt generate password

# Generate a 32-character password without symbols
utility-belt generate password --length 32 --no-symbols

# Generate 5 passwords, one per line
utility-belt generate password --count 5

# Copy a 20-character password directly to clipboard (macOS)
utility-belt generate password --length 20 | pbcopy

# Generate a 6-digit PIN
utility-belt generate password --length 6 --no-symbols --no-uppercase --no-lowercase
```

#### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--length` | `int` | `16` | Exact length of each password (must be > 0) |
| `--count` | `int` | `1` | Number of passwords to generate (one per line, must be > 0) |
| `--no-symbols` | `bool` | `false` | Exclude symbols (`!@#$%^&*()-_=+[]{}|;:,.<>?/~`) |
| `--no-numbers` | `bool` | `false` | Exclude digits (`0-9`) |
| `--no-uppercase` | `bool` | `false` | Exclude uppercase letters (`A-Z`) |
| `--no-lowercase` | `bool` | `false` | Exclude lowercase letters (`a-z`) |

---

### UUID Generator

```bash
# Generate a standard UUID v4
utility-belt generate uuid

# Generate a UUID v7 (RFC 9562 timestamp-based)
utility-belt generate uuid --v7

# Generate 3 UUIDs, one per line
utility-belt generate uuid --count 3

# Save 10 UUIDs to a file
utility-belt generate uuid --count 10 > uuids.txt
```

#### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--v7` | `bool` | `false` | Generate UUID version 7 (timestamp-based) instead of version 4 |
| `--count` | `int` | `1` | Number of UUIDs to generate (one per line, must be > 0) |

---

### Base64 Encoder

```bash
# Encode text via stdin pipe
echo -n "hello" | utility-belt encode base64

# Encode text via --string flag
utility-belt encode base64 --string "hello"

# Copy encoded output to clipboard (macOS)
utility-belt encode base64 --string "my-secret" | pbcopy
```

#### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--string` | `string` | `""` | Text to encode into Base64 (alternative to stdin) |

---

### Base64 Decoder

```bash
# Decode via stdin pipe
echo -n "aGVsbG8=" | utility-belt decode base64

# Decode via --string flag
utility-belt decode base64 --string "aGVsbG8="
```

Invalid or malformed Base64 input produces an error message on `stderr` and exits with code `1`. No partial output is written to `stdout`.

#### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--string` | `string` | `""` | Base64-encoded string to decode (alternative to stdin) |

---

### JWT Decoder

Decodes the Header and Payload of a JSON Web Token (JWS compact format) into readable JSON.

> **Note:** This command performs structural decoding only. It does **not** verify the cryptographic signature.

```bash
# Decode JWT via stdin pipe
echo -n "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" | utility-belt decode jwt

# Decode JWT via --string flag
utility-belt decode jwt --string "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Pipe output to jq for further processing
utility-belt decode jwt --string "<token>" | jq '.payload.exp'
```

Output is a pretty-printed JSON object with `header` and `payload` fields:

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "1234567890",
    "name": "John Doe"
  }
}
```

If the token contains an `exp` claim that has already passed, a warning is printed to `stderr` while the JSON is still written to `stdout` and the command exits with code `0`:

```
warning: token is expired
```

#### Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--string` | `string` | `""` | JWT token string to decode (alternative to stdin) |

---

## Running Tests

```bash
go test -v -race ./...
```
