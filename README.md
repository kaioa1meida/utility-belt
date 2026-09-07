# utility-belt

A fast, modular, offline, and stateless CLI tool designed for developers' everyday tasks.

## Features

- 🔐 **Password Generator:** Cryptographically secure passwords with configurable length and character classes (`crypto/rand`).
- 🆔 **UUID Generator:** Canonical UUID generation supporting UUID v4 (random) and UUID v7 (timestamp-based per RFC 9562).
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

## Running Tests

```bash
go test -v -race ./...
```
